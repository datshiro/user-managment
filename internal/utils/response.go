package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseObject struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func ResponseWithJSON(c *gin.Context, data interface{}) {
	if gin.IsDebugging() {
		c.IndentedJSON(http.StatusOK, ResponseObject{Success: true, Data: data})
	} else {
		c.JSON(http.StatusOK, ResponseObject{Success: true, Data: data})
	}
}

