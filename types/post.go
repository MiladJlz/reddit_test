package types

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

const (
	minTitleLen       = 6
	minDescriptionLen = 10
)

type Post struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      uuid.UUID `gorm:"type:uuid" json:"userID"`
	Votes       []Vote    `gorm:"foreignKey:PostID;references:ID" json:"-"`
	VoteCount   int       `gorm:"default:0" json:"voteCount"`
	CreatedAt   time.Time
}
type CreatePostParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (params CreatePostParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.Title) < minTitleLen {
		errors["title"] = fmt.Sprintf("title length should be at least %d characters", minTitleLen)
	}
	if len(params.Description) < minDescriptionLen {
		errors["description"] = fmt.Sprintf("description length should be at least %d characters", minDescriptionLen)
	}

	return errors
}
func NewPostFromParams(params CreatePostParams) Post {
	return Post{
		Title:       params.Title,
		Description: params.Description,
	}
}
