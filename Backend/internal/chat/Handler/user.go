package user

import (
	"fmt"
	// "log/slog"

	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/storage/postgres"
	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/types"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	DB *postgres.Postgres
}

type response struct {
	Status  bool   `json:"Status"`
	Message string `json:"Message"`
	Data    any    `json:"Data,omitempty"`
}

func (h *Handler) Registerhandler(c *fiber.Ctx) error {
	user := new(types.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{
			Status:  false,
			Message: "Invalid request body",
		})
	}
	// fmt.Println(user)

	// Make sure these field names match your types.User struct!
	if user.Username == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response{
			Status:  false,
			Message: "Name and Password are required",
		})
	}

	res := register(h.DB, user)

	if !res.Status {
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func register(db *postgres.Postgres, user *types.User) *response {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return &response{Status: false, Message: "Error processing password"}
	}
	user.Password = string(hashPassword)

	dbUser, err := db.RegisterUser(user)
	if err != nil {
		return &response{Status: false, Message: "Error registering user in database"}
	}

	return &response{
		Status:  true,
		Message: "User registered successfully",
		Data:    dbUser,
	}
}