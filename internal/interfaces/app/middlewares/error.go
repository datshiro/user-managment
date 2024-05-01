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
      log.Println("ErrorHandlerMiddleware", errs)
			err, ok := errs[0].Err.(consts.CakeError)
			if ok {
        // Log root error
        log.Println(err.Details())
        log.Println(err.Error())
				c.JSON(http.StatusBadRequest, err)
        return
			}
      log.Printf("Unknown err %+v \n", err)
      c.JSON(http.StatusInternalServerError, err)
		}
	}
}
