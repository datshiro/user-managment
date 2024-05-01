package register

import (
	"app/internal/handler"
	"app/internal/utils"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type Request interface {
  handler.Request 
  Marshal() ([]byte, error)
}

func NewRequest() Request {
	return &request{}
}

type request struct {
	Username string `form:"username" binding:"required" json:"username"`
	Email    string `form:"email" binding:"required" json:"email"`
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

