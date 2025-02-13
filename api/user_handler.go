package api

import (
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
		return c.Status(fiber.StatusInternalServerError).JSON(ErrFailedCreatingUUID("user"))
	}
	user.ID = userID
	res, err := h.userStore.InsertUser(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrFailedInsertingData("user"))
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}
func ValidateCreateUser(c fiber.Ctx) error {
	userParam := types.CreateUserParam{}
	if err := c.Bind().Body(&userParam); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrInvalidRequestBody())

	}
	if err := userParam.Validate(); len(err) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	user := types.NewUserFromParams(userParam)
	c.Locals("user", user)
	return c.Next()
}
