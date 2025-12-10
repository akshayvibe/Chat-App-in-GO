package chat

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newLine = []byte{'\n'}
	space   = []byte{' '}
)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

//we will create a constructor of the client struct

func (c *Client) readPump() {
	//the cleanup before exiting the function
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	//setting the rules
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add((pongWait)))
	c.Conn.SetPingHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	//reading the message
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newLine, space, -1))
		c.Hub.Broadcast <- message
	}
}
func (c *Client) writePump() {
	// PART 1: The Heartbeat Clock (Ghadi)
	ticker := time.NewTicker(pingPeriod)

	// PART 2: Cleanup (Jaate waqt ghadi band karna)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	// PART 3: The Traffic Police (Decision Maker)
	for {
		select {

		// CASE A: Hub ne Message Bheja
		case message, ok := <-c.Send:
			// Deadline set karo (10 sec mein bhej dena warna error)
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			// Agar Hub ne channel band kar diya (Server shutdown)
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Message likhna shuru karo
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// OPTIMIZATION: Agar queue mein aur messages hain,
			// toh unhe bhi isi lifafe (packet) mein daal do.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newLine)
				w.Write(<-c.Send)
			}

			// Lifaafa band karke bhej do
			if err := w.Close(); err != nil {
				return
			}

		// CASE B: Ghadi ki Ghanti Baji (Ping Time)
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			// Chupke se ek 'Ping' bhejo user ko
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	go client.readPump()
	go client.writePump()
}
