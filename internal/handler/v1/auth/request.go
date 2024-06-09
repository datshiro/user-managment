package auth

import (
	"app/internal/models"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type LoginRequest struct {
	Email    string `json:"email" `
	Password string `json:"password" `
}

func ValidateRegisterRequest(user *models.User) error {
	if err := validation.Validate(user.Email, validation.Required); err != nil {
		return errors.New("email must not be empty")
	}
	if err := validation.Validate(user.Email, is.Email); err != nil {
		return errors.New("email is in correct format")
	}
	if err := validation.Validate(user.Password, validation.Required); err != nil {
		return errors.New("password must not be empty")
	}
	return nil
}
