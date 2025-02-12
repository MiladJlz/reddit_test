package types

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Vote struct {
	UserID uuid.UUID `gorm:"primaryKey;type:uuid" json:"userID" validate:"uuid"`
	PostID uuid.UUID `gorm:"primaryKey;type:uuid;index" json:"postID" validate:"required,uuid"`
	Value  int       `json:"value" validate:"required,voteValue"`
}

func (v *Vote) Validate() error {
	validate := validator.New()
	err := validate.RegisterValidation("voteValue", validateVoteValue)
	if err != nil {
		return err
	}
	return validate.Struct(v)

}
func validateVoteValue(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(int)
	if !ok {
		return false
	}
	return val == -1 || val == 1
}
