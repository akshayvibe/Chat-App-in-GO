package handler

import (
	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/storage/postgres"
	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/types"
	"github.com/gofiber/fiber/v2"
)

type MessageHandler struct {
	DB *postgres.Postgres
}
type MessageResponse struct{
	Message []types.Message `json:"Message"`
}
func (m *MessageHandler) GetMessages(c *fiber.Ctx) error {
    // 1. Get RoomID from the URL parameter (e.g., /messages/:roomId)
    roomIdStr := c.Params("roomId")
    
    // 2. Define a slice to hold the messages
    var messages []types.Message

    // 3. Query Postgres using GORM
    // We use .Preload("User") to automatically join the User table and get sender details
    // We use .Order("created_at asc") so the chat history is in chronological order
    result := m.DB.Db.Preload("User").Where("room_id = ?", roomIdStr).Order("created_at asc").Find(&messages)

    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Could not fetch messages",
        })
    }

    // 4. Return the list of messages as JSON
    return c.Status(fiber.StatusOK).JSON(messages)
}