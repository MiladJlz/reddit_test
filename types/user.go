package types

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	Username string    `json:"username" validate:"required,gte=6"`
	Vote     []Vote    `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)

}
