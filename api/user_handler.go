package api

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/miladjlz/reddit_test/db"
	"github.com/miladjlz/reddit_test/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{userStore: userStore}

}

func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	user := c.Locals("user").(types.User)
	userID, err := uuid.NewUUID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf(" %v", err)})
	}
	user.ID = userID
	res, err := h.userStore.InsertUser(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}
func ValidateCreateUser(c fiber.Ctx) error {
	user := types.User{}
	if err := c.Bind().Body(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"invalid request body": fmt.Sprintf("%v", err)})

	}
	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validation errors": fmt.Sprintf("%v", err)})

	}
	c.Locals("user", user)
	return c.Next()
}
