package middlewares

import (
	"app/internal/common/helper"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

    for _, err := range c.Errors {
			apiErr, ok := err.Err.(helper.APIError)
			if ok {
        c.AbortWithStatusJSON(apiErr.StatusCode, apiErr)
        return
			} else {
        c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any {
          "status_code": http.StatusInternalServerError,
          "message": "internal server error",
        })
      }
      slog.Error("HTTP API error", "err", err, "path", c.Request.URL.Path)
    }
	}
}
