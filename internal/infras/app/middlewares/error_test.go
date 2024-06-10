package middlewares

import (
	"app/internal/common/helper"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

type handler struct {
}

func (h handler) Handle(c *gin.Context) {

}

func TestErrorHandlerMiddleware(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  error
		wantResp string
		wantCode int
	}{
		{
			"Test handler return invalid request",
			errors.New(""),
			`{"message":"internal server error","statusCode":500}`,
			500,
		},
		{
			"Test handler return APIError",
      helper.NewAPIError(http.StatusBadRequest, errors.New("invalid request JSON data")),
      `{"status_code":400,"message":"invalid request JSON data"}`,
			400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := gin.New()
			rec := httptest.NewRecorder()
			h := func(c *gin.Context) { c.Error(tt.wantErr) }

			app.Use(ErrorHandlerMiddleware())
			app.GET("/", h)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			app.ServeHTTP(rec, req)

			if got := rec.Code; !reflect.DeepEqual(got, tt.wantCode) {
				t.Errorf("ErrorHandlerMiddleware() = %v, want %v", got, tt.wantCode)
			}

			if got := rec.Body.String(); !reflect.DeepEqual(got, tt.wantResp) {
				t.Errorf("ErrorHandlerMiddleware() = %v, want %v", got, tt.wantResp)
			}
		})
	}
}
