package login

import (
	"app/internal/consts"
	"app/internal/handler"
	"app/internal/interfaces/usecases/user"
	"app/internal/models"
	"app/internal/utils"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewHandler(userUC user.UserUsecase) handler.Handler {
	return loginHandler{NewRequest: NewRequest, UserUsecase: userUC}
}

type loginHandler struct {
	NewRequest  func() Request
	UserUsecase user.UserUsecase
}

func (h loginHandler) Handle(c *gin.Context) {
	req := h.NewRequest()
	if err := req.Bind(c); err != nil {
		c.Error(err)
		return
	}

	if err := req.Validate(); err != nil {
		c.Error(err)
		return
	}

	var (
		user *models.User
		err  error
	)

	switch req.LoginType() {
	case consts.UserNameLoginType:
		user, err = h.UserUsecase.LoginWithUsername(c.Request.Context(), req.GetAccount(), req.GetPassword())
	case consts.EmailLoginType:
		user, err = h.UserUsecase.LoginWithEmail(c.Request.Context(), req.GetAccount(), req.GetPassword())
	case consts.PhoneNumberLoginType:
		user, err = h.UserUsecase.LoginWithPhone(c.Request.Context(), req.GetAccount(), req.GetPassword())
	default:
		c.Error(consts.NewError(fmt.Errorf("provided login does not support")).
			WithRootCause(err).
			WithTag("Method", "LoginUser").
			WithTag("data", req))
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(consts.ErrDataNotFound.WithHttpCode(http.StatusNotFound).WithTag("Method", "LoginUser").WithTag("data", req))
			return
		}
		c.Error(consts.ErrLoginFailure.WithRootCause(err).WithTag("Method", "LoginUser").WithTag("data", req))
		return
	}

	utils.ResponseWithJSON(c, toResponse(user))
}

func toResponse(user *models.User) response {
  return response{
    Fullname   : user.Fullname,
    PhoneNumber: user.PhoneNumber,
    Email      : user.Email,
    Username   : user.Username,
    Birthday   : user.Birthday,
    LatestLogin: user.LatestLogin,
  }
}

type response struct {
	Fullname    string    `json:"fullname"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Birthday    time.Time `json:"birthday"`
	LatestLogin time.Time `json:"latest_login"`
}
