package middlewares

import (
	"app/internal/consts"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		errs := c.Errors
		if len(errs) > 0 {
			err, ok := errs[0].Err.(consts.CakeError)
			if ok {
				// Log root error
        log.Printf("error details: %v",err.Details())
				responseCakeError(c, err)
				return
			}
			log.Printf("Unknown err %+v \n", err)
			c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err})
		}
	}
}

func responseCakeError(c *gin.Context, err consts.CakeError) {
	c.JSON(err.GetCode(), map[string]interface{}{"error": err.Error()})
}
