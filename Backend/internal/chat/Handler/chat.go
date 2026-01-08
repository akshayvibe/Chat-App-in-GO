package handler

import (
	"math/rand"
	"time"

	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/storage/postgres"
	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/types"
	"github.com/gofiber/fiber/v2"
)

type Chathandler struct {
	DB *postgres.Postgres
}
type ChatRoomResponse struct {
	ID   uint   `json:"Id"`
	Name string `json:"RoomName"`
	Code string `json:"Code"`
}

func (ch *Chathandler) CreateChatRoom(c *fiber.Ctx) error {
	chatRoom := &types.Room{}
	if err := c.BodyParser(chatRoom); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  false,
			Message: "Invalid Request Body",
		})
	}
	if chatRoom.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  false,
			Message: "Name is required",
		})
	}
	res := registerRoom(ch.DB, chatRoom)
	if !res.Status {
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

func registerRoom(db *postgres.Postgres, chatRoom *types.Room) *Response {
	Code := String()
	chatRoom.RoomCode = Code
	dbRoom, err := db.RegisterRoom(chatRoom)
	if err != nil {
		return &Response{Status: false, Message: "Error registering ChatRoom in databases"}
	}
	return &Response{
		Status:  true,
		Message: "Chat room registered successfully",
		Data: &ChatRoomResponse{
			ID:   dbRoom.ID,
			Name: dbRoom.Name,
			Code: dbRoom.RoomCode,
		},
	}
}
func (ch *Chathandler) JoinRoom(c *fiber.Ctx) error {
	type Request struct {
		Code   string `json:"code"`
		UserID uint   `json:"user_id"`
	}
	var req Request

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid body"})
	}
	res := joinRoom(ch.DB, req.Code, req.UserID)
	return c.Status(201).JSON(res)
}
func joinRoom(db *postgres.Postgres, code string, userID uint) *Response {
	room, err := db.GetRoom(code)
	if err != nil {
		return &Response{Status: false, Message: "Room not found"}
	}
	member := types.RoomMember{
		RoomID: room.ID,
		UserID: userID,
		Role:   "member",
	}
	existing, _ := db.CheckExistingMembers(userID, room.ID)
	if existing != nil {
		return &Response{Status: false, Message: "User already exists in the room"}
	}
	if err := db.Db.Create(&member).Error; err != nil {
		return &Response{Status: false, Message: "Could not join room"}
	}

	return &Response{Status: true, Message: "Joined successfully"}
}

// CREATING RANDOM STRING CODE
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String() string {
	return StringWithCharset(5, charset)
}
