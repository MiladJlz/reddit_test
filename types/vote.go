package types

import (
	"fmt"
	"github.com/google/uuid"
)

type Vote struct {
	UserID uuid.UUID `gorm:"primaryKey;type:uuid" json:"userID"`
	PostID uuid.UUID `gorm:"primaryKey;type:uuid;index" json:"postID"`
	Value  int       `json:"value"`
}

type CreateVoteParams struct {
	PostID string `json:"postID"`
	Value  int    `json:"value"`
}

func (params CreateVoteParams) Validate() map[string]string {
	errors := map[string]string{}
	if !VoteValueValidation(params.Value) {
		errors["value"] = fmt.Sprintf("value must be -1 or 1")
	}
	if err := PostIDValidation(params.PostID); err != nil {
		errors["postID"] = fmt.Sprintf("postID must be a UUID")
	}
	return errors
}
func PostIDValidation(id string) error {
	_, err := uuid.Parse(id)
	return err
}

func VoteValueValidation(value int) bool {
	if value == 1 || value == -1 {
		return true
	}
	return false
}
func NewVoteFromParams(params CreateVoteParams) Vote {
	postID, _ := uuid.Parse(params.PostID)
	return Vote{
		PostID: postID,
		Value:  params.Value,
	}
}
