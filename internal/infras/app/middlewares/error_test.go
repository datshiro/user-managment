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

func TestErrorHandlerMiddleware(t *testing.T) {
	tests := []struct {
		name     string
		wantErr  error
		wantResp string
		wantCode int
	}{
		{
			"Case return invalid request",
			errors.New("this is unknown error"),
			`{"message":"internal server error","status_code":500}`,
			http.StatusInternalServerError,
		},
		{
			"Case return APIError",
			helper.NewAPIError(http.StatusBadRequest, errors.New("invalid request JSON data")),
			`{"status_code":400,"message":"invalid request JSON data"}`,
			http.StatusBadRequest,
		},
		{
			"Case return ValidationErrors",
			helper.InvalidValidationErrors(map[string]error{
				"email":    errors.New("email must not be empty"),
				"password": errors.New("password must not be empty"),
			}),
			`{"status_code":422,"message":{"email":"email must not be empty","password":"password must not be empty"}}`,
			http.StatusUnprocessableEntity,
		},
    {
      "Case return bad request",
      helper.InvalidJSON(),
      `{"status_code":400,"message":"invalid request JSON data"}`,
      http.StatusBadRequest,
    },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
      gin.SetMode(gin.TestMode)
			app := gin.New()
			rec := httptest.NewRecorder()
			app.Use(ErrorHandlerMiddleware())
			h := func(c *gin.Context) { c.Error(tt.wantErr) }

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
