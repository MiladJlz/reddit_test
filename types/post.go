package types

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type Post struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid" json:"id" validate:"uuid"`
	Title       string    `json:"title" validate:"required,gte=6"`
	Description string    `json:"description" validate:"required,gte=10"`
	UserID      uuid.UUID `gorm:"type:uuid" json:"userID"`
	Votes       []Vote    `gorm:"foreignKey:PostID;references:ID" json:"-"`
	VoteCount   int       `gorm:"default:0" json:"voteCount"`
	CreatedAt   time.Time
}

func (p *Post) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
