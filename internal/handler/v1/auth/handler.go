package auth

import (
	"app/internal/common/http/response"
	"app/internal/interfaces/usecases/user"
	"app/internal/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gorm.io/gorm"
)

type AuthHandler interface {
	HandleLogin(*gin.Context)
	HandleRegister(*gin.Context)
}

type handler struct {
	uc user.UserUsecase
}

func NewHandler(uc user.UserUsecase) AuthHandler {
	return handler{uc}
}

// Handlers
func (h handler) HandleLogin(c *gin.Context) {
	var req = &LoginRequest{}

	if err := c.BindJSON(req); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if err := validation.Validate(req.Email, validation.Required, is.Email); err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("email must not be empty"))
		return
	}

	if err := validation.Validate(req.Password, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("password must not be empty"))
		return
	}

	user, err := h.uc.LoginWithEmail(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.New(c).Error(http.StatusNotFound, err)
			return
		}
		response.New(c).Error(http.StatusInternalServerError, fmt.Errorf("something went wrong, please try again!; %v", err))
		return
	}

	// TODO: generate token

	response.New(c).JSON(user)

}

func (h handler) HandleRegister(c *gin.Context) {
	var req = &models.User{}

	if err := c.BindJSON(req); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if err := ValidateRegisterRequest(req); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
    return
	}

	user, err := h.uc.RegisterUser(c.Request.Context(), req)
	if err != nil {
		response.New(c).Error(http.StatusConflict, err)
		return
	}

	response.New(c, response.WithCode(http.StatusCreated)).JSON(user)
}
