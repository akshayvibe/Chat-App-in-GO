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
	var roomId uint
	if err := c.BodyParser(&roomId); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  false,
			Message: "Invalid Request Body",
		})
	}
	res := getRoomMess(m.DB, roomId)
	if !res.Status {
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}
	return c.Status(fiber.StatusCreated).JSON(res)
	// return c.JSON(messages)
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
func (m *MessageHandler)GetPrivateMessage(c *fiber.Ctx)error{

	var UserA uint
	if err := c.BodyParser(&UserA); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  false,
			Message: "Invalid Request Body",
		})
	}
	var UserB uint
	if err := c.BodyParser(&UserB); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  false,
			Message: "Invalid Request Body",
		})
	}
	res:=getPrivateMessage(m.DB,UserA,UserB)
	if !res.Status {
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}

func getPrivateMessage(db *postgres.Postgres,userA,userB uint)*Response{
	privateMessage,err := db.getPrivateMessages(userA,userB)
	if err != nil {
		return &Response{
			Status:  false,
			Message: "Could not retrieve rooms",
		}
	}
	return &Response{
		Status:  true,
		Message: "Room Message retrieved successfully",
		Data:    privateMessage,
	}
}
