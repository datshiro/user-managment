package auth

import (
	"app/internal/interfaces/repositories/models"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func (req LoginRequest) Validate() map[string]error {
  errs := make(map[string]error)
	if err := validation.Validate(req.Email, validation.Required, is.Email); err != nil {
		errs["email"] = errors.New("email must not be empty")
	}

	if err := validation.Validate(req.Password, validation.Required); err != nil {
		errs["password"] = errors.New("email must not be empty")
	}
  if len(errs) > 0 {
    return errs
  }
  return nil
}

func ValidateRegisterRequest(user *models.User) map[string]error {
  errs := make(map[string]error) 
	if err := validation.Validate(user.Email, validation.Required); err != nil {
		errs["message"] = errors.New("email must not be empty")
	}
	if err := validation.Validate(user.Email, is.Email); err != nil {
		errs["email_format"] = errors.New("email is in correct format")
	}
	if err := validation.Validate(user.Password, validation.Required); err != nil {
		errs["password"] = errors.New("password must not be empty")
	}
  if len(errs) > 0 {
    return errs
  }
	return nil
}
