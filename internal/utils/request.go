package utils

import (
	"app/internal/consts"

	"github.com/gin-gonic/gin"
)


// Wrap function for binding request and handle error
// handlerTag should be the function call 
func BindRequest(c *gin.Context, request interface{}, handlerTag string) error {
	if err := c.Bind(&request); err != nil {
		err = consts.ErrInvalidRequest.WithRootCause(err).WithTag("handler",handlerTag)
		return err
	}
	return nil
}
