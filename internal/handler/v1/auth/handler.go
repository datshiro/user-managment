package auth

import (
	"app/internal/common/helper"
	"app/internal/common/http/response"
	"app/internal/interfaces/repositories/models"
	"app/internal/interfaces/usecases/user"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
	var req = LoginRequest{}

	if err := c.BindJSON(&req); err != nil {
    c.Error(helper.InvalidJSON())
		return
	}

  if errs := req.Validate() ; errs != nil {
    c.Error(helper.InvalidValidationErrors(errs))
    return
  }

	user, err := h.uc.LoginWithEmail(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
      c.Error(helper.NewAPIError(http.StatusNotFound, errors.New("user is not exist")))
			return
		}
    c.Error(err)
		return
	}

	// TODO: generate token

  c.JSON(http.StatusOK, map[string]any{
    "email": user.Email,
    "birthday": user.Birthday,
    "phone": user.PhoneNumber,
  })
}

func (h handler) HandleRegister(c *gin.Context) {
	var req = &models.User{}

	if err := c.BindJSON(req); err != nil {
    c.Error(helper.InvalidJSON().WithDetails(err))
		return
	}

	if err := ValidateRegisterRequest(req); err != nil {
    c.Error(helper.InvalidValidationErrors(err))
    return
	}

	user, err := h.uc.RegisterUser(c.Request.Context(), req)
	if err != nil {
		response.New(c).Error(http.StatusConflict, err)
		return
	}

	response.New(c, response.WithCode(http.StatusCreated)).JSON(user)
}
