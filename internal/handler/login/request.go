package login

import (
	"app/internal/consts"
	"app/internal/handler"
	"app/internal/utils"

	"github.com/gin-gonic/gin"
)

type Request interface {
	handler.Request
	GetAccount() string
	GetPassword() string
  LoginType() consts.LoginType
}

func NewRequest() Request {
	return &request{}
}

type request struct {
	Username    string `form:"username"  json:"username"`
	Email       string `form:"email"  json:"email"`
	PhoneNumber string `form:"phone_number" json:"phone_number"`
	Password    string `form:"password" json:"password"`
	account     string // Account value for authenticating
	loginType   consts.LoginType // Username | Email | PhoneNumber
}

func (r *request) Bind(c *gin.Context) error {
	if err := utils.BindRequest(c, r, "LoginHandler"); err != nil {
		return err
	}
	if r.Username != "" {
		r.account = r.Username
    r.loginType = consts.UserNameLoginType
	}
	if r.Email != "" {
		r.account = r.Email
    r.loginType = consts.EmailLoginType
	}
	if r.PhoneNumber != "" {
		r.account = r.PhoneNumber
    r.loginType = consts.PhoneNumberLoginType
	}
	return nil
}

func (r *request) Validate() error {
	if r.account == "" && r.loginType == 0 {
		return consts.ErrMissingCredentialInfo
	}
	if r.Password == "" {
    return consts.ErrPasswordCannotBeEmpty
	}
  return nil
}

func (r *request) GetAccount() string {
	return r.account
}

func (r *request) GetPassword() string {
	return r.Password
}

func (r *request) LoginType() consts.LoginType {
  return r.loginType
}
