package register

import (
	"app/internal/handler"
	"app/internal/interfaces/usecases/user"
	"app/internal/utils"

	"github.com/gin-gonic/gin"
)

func NewHandler(userUC user.UserUsecase) handler.Handler {
	return registerHandler{NewRequest: NewRequest, UserUsecase: userUC}
}

type registerHandler struct {
	NewRequest func() Request
  UserUsecase user.UserUsecase
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
  h.UserUsecase.RegisterUser()

	utils.ResponseWithJSON(c, req)
}
