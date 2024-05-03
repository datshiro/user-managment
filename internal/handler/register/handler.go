package register

import (
	"app/internal/consts"
	"app/internal/handler"
	"app/internal/interfaces/usecases/user"
	"app/internal/utils"
	"fmt"

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

  userData , err := req.toData()
  if err !=nil {
		c.Error(consts.ErrInvalidRequest.WithRootCause(err).WithTag("Method", "ConvertRequestData").WithTag("request", fmt.Sprintf("%+v", req)))
    return 
  }

  user, err := h.UserUsecase.RegisterUser(c.Request.Context(), userData)
  if err !=nil {
		c.Error(consts.ErrCreateFailure.WithRootCause(err).WithTag("Method", "RegisterUser").WithTag("data", userData))
    return 
  }

	utils.ResponseWithJSON(c, user)
}
