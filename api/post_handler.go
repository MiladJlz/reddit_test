package api

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/miladjlz/reddit_test/db"
	"github.com/miladjlz/reddit_test/types"
	"strings"
)

type PostHandler struct {
	postStore db.PostStore
}

func NewPostHandler(postStore db.PostStore) *PostHandler {
	return &PostHandler{postStore: postStore}

}

func (h *PostHandler) CreatePost(c fiber.Ctx) error {
	post := c.Locals("post").(types.Post)
	postID, err := uuid.NewUUID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed creating uuid for post"})
	}
	post.ID = postID
	res, err := h.postStore.InsertPost(&post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed inserting post to DB"})
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}
func ValidateCreatePost(c fiber.Ctx) error {
	post := types.Post{}
	if err := c.Bind().Body(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	param := c.Params("id")
	userID, err := uuid.Parse(param)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"invalid userID": " must be a UUID"})
	}
	if err := post.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validation errors": fmt.Sprintf(" %v", err)})
	}
	post.UserID = userID
	c.Locals("post", post)
	return c.Next()
}

func (h *PostHandler) GetPosts(c fiber.Ctx) error {

	posts, err := h.postStore.GetPosts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf(": %v", err)})
	}
	return c.Status(fiber.StatusOK).JSON(posts)
}

func (h *PostHandler) GetPostsByVote(c fiber.Ctx) error {

	posts, err := h.postStore.GetPostsByVote()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf(": %v", err)})
	}
	return c.Status(fiber.StatusOK).JSON(posts)
}

func (h *PostHandler) GetPostsByOrder(c fiber.Ctx) error {
	order := c.Locals("order").(string)
	posts, err := h.postStore.GetPostsByDate(order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf(": %v", err)})
	}
	return c.Status(fiber.StatusOK).JSON(posts)
}
func ValidateGetPostsByOrder(c fiber.Ctx) error {

	order := c.Query("order")
	orderToUpper := strings.ToUpper(order)
	if orderToUpper != "ASC" && orderToUpper != "DESC" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid order parameter"})
	}
	c.Locals("order", orderToUpper)
	return c.Next()
}
