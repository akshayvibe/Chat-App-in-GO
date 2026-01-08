package handler

import (
	// "fmt"
	// "log/slog"

	"errors"

	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/storage/postgres"
	"github.com/akshayjha21/Chat-App-in-GO/Backend/internal/types"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *postgres.Postgres
}
type UserResponse struct {
	Id       uint   `json:"Id"`
	Name string `json:"Username"`
}
func (h *UserHandler) Registerhandler(c *fiber.Ctx) error {
	user := new(types.User)

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  false,
			Message: "Invalid request body",
		})
	}
	// fmt.Println(user)

	// Make sure these field names match your types.User struct!
	if user.Username == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
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
func (h *UserHandler) LoginHandler(c *fiber.Ctx) error {
	user := &types.User{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  false,
			Message: "Invalid request body",
		})
	}
	if user.Username == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Status:  false,
			Message: "Name and Password are required",
		})
	}
	res := login(h.DB, user)
	if !res.Status {
		// Use 401 for bad credentials, 500 for actual DB crashes
		if res.Message == "Invalid credentials" {
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
func login(db *postgres.Postgres, user *types.User) *Response {
	u, err := db.GetUser(&types.User{Username: user.Username})
	if err != nil {
		// Specifically check if the record just wasn't found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &Response{Status: false, Message: "Invalid credentials"}
		}
		// Handle actual UserResponse base connection errors
		return &Response{Status: false, Message: "Internal server error"}
	}
	if !checkPassword(user.Password, u.Password) {
		return &Response{Status: false, Message: "Invalid password"}
	}
	return &Response{
		Status:  true,
		Message: "User login succesfully",
		// Using map[string]any if your Response Data field supports it
		Data: &UserResponse{
			Id:   u.ID,
			Name: u.Username},
	}
}
func register(db *postgres.Postgres, user *types.User) *Response {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return &Response{Status: false, Message: "Error processing password"}
	}
	user.Password = string(hashPassword)

	dbUser, err := db.RegisterUser(user)
	if err != nil {
		return &Response{Status: false, Message: "Error registering user in Database"}
	}

	return &Response{
		Status:  true,
		Message: "User registered succesfully",
		Data: &UserResponse{
			Id:   dbUser.ID,
			Name: dbUser.Username,
		},
	}
}

func checkPassword(plainPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

//TODO -  JWT Token based authentication
//TODO - RBAC  