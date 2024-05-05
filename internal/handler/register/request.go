package register

import (
	"app/internal/consts"
	"app/internal/handler"
	"app/internal/interfaces/entities"
	"app/internal/utils"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

type Request interface {
	handler.Request
	Marshal() ([]byte, error)
	toData() (entities.UserData, error)
}

func NewRequest() Request {
	return &request{}
}

type request struct {
	Username    string           `form:"username"  json:"username"`
	Email       string           `form:"email"  json:"email"`
	Fullname    string           `form:"fullname" json:"fullname"`
	PhoneNumber string           `form:"phone_number" json:"phone_number"`
	Password    string           `form:"password" json:"password"`
	Birthday    time.Time        `form:"birthday" json:"birthday"`
	LatestLogin time.Time        `form:"latest_login" json:"latest_login"`
	account     string           // Account value for authenticating
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

func (r *request) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *request) toData() (entities.UserData, error) {
	b, err := r.Marshal()
	if err != nil {
		return entities.UserData{}, err
	}

	data := entities.UserData{}
	if err := json.Unmarshal(b, &data); err != nil {
		return entities.UserData{}, err
	}
	return data, nil
}
