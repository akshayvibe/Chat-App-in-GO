package handler

import (
	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/storage/postgres"
	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/types"
	"github.com/gofiber/fiber/v2"
)

type MessageHandler struct {
	DB *postgres.Postgres
}
type MessageResponse struct {
	Message []types.Message `json:"Message"`
}

func (m *MessageHandler) GetRoomMessages(c *fiber.Ctx) error {
    roomId, err := c.ParamsInt("roomId")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(Response{
            Status:  false,
            Message: "Invalid Room ID",
        })
    }
    res := getRoomMess(m.DB, uint(roomId))
    if !res.Status {
        return c.Status(fiber.StatusInternalServerError).JSON(res)
    }
    return c.Status(fiber.StatusOK).JSON(res)
}
func getRoomMess(db *postgres.Postgres, roomId uint) *Response {
	room_msg, err := db.GetRoomMessages(roomId)
	if err != nil {
		return &Response{
			Status:  false,
			Message: "Could not retrieve rooms",
		}
	}
	return &Response{
		Status:  true,
		Message: "Room Message retrieved successfully",
		Data:    room_msg,
	}
}
func (m *MessageHandler) GetPrivateMessage(c *fiber.Ctx) error {
    userA, errA := c.ParamsInt("userA")
    userB, errB := c.ParamsInt("userB")

    if errA != nil || errB != nil {
        return c.Status(fiber.StatusBadRequest).JSON(Response{
            Status:  false,
            Message: "User IDs must be integers",
        })
    }
    res := getPrivateMessage(m.DB, uint(userA), uint(userB))
    
    if !res.Status {
        return c.Status(fiber.StatusInternalServerError).JSON(res)
    }
    return c.Status(fiber.StatusOK).JSON(res)
}

func getPrivateMessage(db *postgres.Postgres, userA, userB uint) *Response {
    privateMessages, err := db.GetPrivateMessages(userA, userB)
    
    if err != nil {
        return &Response{
            Status:  false,
            Message: "Failed to retrieve conversation history",
        }
    }
    return &Response{
        Status:  true,
        Message: "Conversation retrieved successfully",
        Data:    privateMessages,
    }
}
//ANCHOR - have to add paggination in the message api
