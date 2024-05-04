package login

import (
	"app/internal/consts"
	"app/internal/handler"
	"app/internal/interfaces/usecases/user"
	"app/internal/models"
	"app/internal/utils"
	"errors"
	"net/http"

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
		c.Error(consts.ErrLoginFailure.WithRootCause(err).WithTag("Method", "LoginUser").WithTag("data", req))
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Error(consts.ErrDataNotFound.WithHttpCode(http.StatusNotFound).WithTag("Method", "LoginUser"))
			return
		}
		c.Error(consts.ErrLoginFailure.WithRootCause(err).WithTag("Method", "LoginUser").WithTag("data", req))
		return
	}

	utils.ResponseWithJSON(c, user)
}