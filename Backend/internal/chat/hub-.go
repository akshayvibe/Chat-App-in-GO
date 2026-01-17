package chat

import (
	"log"
	"fmt"
	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/storage/postgres"
)

// import "golang.org/x/text/message"

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	db         *postgres.Postgres
}

func NewHub(db *postgres.Postgres) *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		db:         db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			log.Println("disconnected from", client.Conn.RemoteAddr())
		case message := <-h.Broadcast:
			fmt.Println("new message", message)
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
