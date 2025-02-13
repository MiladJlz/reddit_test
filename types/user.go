package types

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	minUsernameLen = 6
)

type User struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	Username string    `json:"username"`
	Vote     []Vote    `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

type CreateUserParam struct {
	Username string `json:"username"`
}

func (params CreateUserParam) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.Username) < minUsernameLen {
		errors["title"] = fmt.Sprintf("username length should be at least %d characters", minUsernameLen)
	}
	return errors
}
func NewUserFromParams(params CreateUserParam) User {
	return User{
		Username: params.Username,
	}
}
