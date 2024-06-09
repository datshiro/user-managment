
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseObject struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

func JSON(c *gin.Context, data interface{}) {
	if gin.IsDebugging() {
		c.IndentedJSON(http.StatusOK, ResponseObject{Success: true, Data: data})
	} else {
		c.JSON(http.StatusOK, ResponseObject{Success: true, Data: data})
	}
}

func Error(c *gin.Context, code int, err error, data interface{}, messages ...string) {
  message := err.Error()
  for _, msg := range messages {
    message += msg
  }
	c.JSON(code, ResponseObject{Success: false, Data: data, Message: message})
}
