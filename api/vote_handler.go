package api

import (
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
	res, err := h.postStore.HasUserVoted(vote.UserID, vote.PostID)
	if res == false && err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrFailedGettingData("vote"))
	} else if res == false {
		res, err := h.voteStore.InsertVote(&vote)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrFailedInsertingData("vote"))
		}
		return c.Status(fiber.StatusCreated).JSON(res)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(ErrVoteConflict())
	}

}
func ValidateAddVote(c fiber.Ctx) error {
	voteParams := types.CreateVoteParams{}
	if err := c.Bind().Body(&voteParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrInvalidRequestBody())
	}
	param := c.Params("id")
	userID, err := uuid.Parse(param)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrInvalidPathParam())
	}
	if errors := voteParams.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors})
	}
	vote := types.NewVoteFromParams(voteParams)
	vote.UserID = userID
	c.Locals("vote", vote)
	return c.Next()
}
