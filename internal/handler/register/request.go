package register

import (
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
	Username    string    `form:"username"  json:"username"`
	Email       string    `form:"email"  json:"email"`
	Fullname    string    `form:"fullname" json:"fullname"`
	PhoneNumber string    `form:"phone_number" json:"phone_number"`
	Password    string    `form:"password" json:"password"`
	Birthday    time.Time `form:"birthday" json:"birthday"`
	LatestLogin time.Time `form:"latest_login" json:"latest_login"`
}

func (r *request) Bind(c *gin.Context) error {
	return utils.BindRequest(c, r, "registerHandler")
}

func (*request) Validate() error {
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
