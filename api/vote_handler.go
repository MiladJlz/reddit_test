package api

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/miladjlz/reddit_test/db"
	"github.com/miladjlz/reddit_test/types"
)

type VoteHandler struct {
	voteStore db.VoteStore
	postStore db.PostStore
}

func NewVoteHandler(voteStore db.VoteStore, postStore db.PostStore) *VoteHandler {
	return &VoteHandler{voteStore: voteStore, postStore: postStore}

}

func (h *VoteHandler) AddVote(c fiber.Ctx) error {
	vote := c.Locals("vote").(types.Vote)
	b, err := h.postStore.HasUserVoted(vote.UserID, vote.PostID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf(" %v", err)})
	}
	if b == false {
		res, err := h.voteStore.InsertVote(&vote)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf(" %v", err)})
		}
		return c.Status(fiber.StatusCreated).JSON(res)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user already voted"})
	}

}
func ValidateAddVote(c fiber.Ctx) error {
	vote := types.Vote{}
	if err := c.Bind().Body(&vote); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"invalid request body": fmt.Sprintf(" %v", err)})
	}
	if err := vote.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validation errors": fmt.Sprintf(" %v", err)})

	}
	param := c.Params("id")
	userID, err := uuid.Parse(param)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"invalid userID": " must be a UUID"})
	}
	vote.UserID = userID
	c.Locals("vote", vote)
	return c.Next()
}
