package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseObject struct {
  Success bool
  Data interface{}
}

func ResponseWithJSON(c *gin.Context, data interface{})  {
  c.JSON(http.StatusOK, ResponseObject{Success: true, Data: data})
}

