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
      log.Println("errs", errs)
			err, ok := errs[0].Err.(consts.CustomError)
			if ok {
				// Log root error
				// log.Printf("error details: %v", err.Details())
				c.JSON(err.GetCode(), err.Error())
				return
			}
			// log.Printf("Unknown err %+v \n", err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
}
