package register

import (
	"app/internal/handler"
	"app/internal/utils"

	"github.com/gin-gonic/gin"
)

func NewHandler() handler.Handler {
	return registerHandler{NewRequest: NewRequest}
}

type registerHandler struct {
	NewRequest func() Request
}

func (h registerHandler) Handle(c *gin.Context) {
	req := h.NewRequest()
	if err := req.Bind(c); err != nil {
		c.Error(err)
    return 
	}

	if err := req.Validate(); err != nil {
		c.Error(err)
    return 
	}

	utils.ResponseWithJSON(c, req)
	return
}
