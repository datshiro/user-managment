package auth

import (
	"app/internal/interfaces/usecases/user"
	mocks "app/internal/mocks/usecases/user"
	"app/internal/models"
	"app/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(uc user.UserUsecase) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.POST("/register", NewHandler(uc).HandleRegister)
	return r
}

func TestSuccessfulRegisterRoute(t *testing.T) {
	tests := []struct {
		name     string
		args     *models.User
		expected *models.User
	}{
		{
			"success: register with email",
			&models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "",
				Email:       "datshiro@gmail.com",
				Username:    "",
				Password:    "sTr0ngP@ssword",
			},
			&models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "",
				Email:       "datshiro@gmail.com",
				Username:    "",
				Password:    "sTr0ngP@ssword",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := new(mocks.UserUsecase)
			mock.
				On("RegisterUser", context.Background(), tt.args).
				Return(tt.expected, nil)

			router := setupRouter(mock)
			w := httptest.NewRecorder()
			b, err := json.Marshal(tt.args)
			assert.NoError(t, err, "Marshal request body failed")

			req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(b))
			assert.NoError(t, err, "Create test request failed")
			router.ServeHTTP(w, req)

			// Parse success result
			expectedResponse := utils.ResponseObject{Success: true, Data: tt.expected}
			expected, err := json.Marshal(expectedResponse)
			assert.NoError(t, err, "Marshal request body failed")

			assert.Equal(t, 201, w.Code, tt.name)
			assert.Equal(t, bytes.NewBuffer(expected).String(), w.Body.String(), tt.name)
		})
	}
}

func TestFailureRegisterRoute(t *testing.T) {
	tests := []struct {
		name     string
		args     *models.User
		expected error
	}{
		{
			"failed: register with empty email",
			&models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "",
				Email:       "datshiro@gmail.com",
				Username:    "",
				Password:    "sTr0ngP@ssword",
			},
      errors.New("email must not be empty"),
		},
		{
			"failed: register with empty password",
			&models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "",
				Email:       "",
				Username:    "",
				Password:    "sTr0ngP@ssword",
			},
      nil,
		},
		{
			"failed: existed user",
			&models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "123456790",
				Email:       "datshiro@gmail.com",
				Username:    "datshiro",
				Password:    "",
			},
      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := new(mocks.UserUsecase)
			mock.
				On("RegisterUser", context.Background(), tt.args).
				Return(nil, tt.expected)

			router := setupRouter(mock)
			w := httptest.NewRecorder()
			b, _ := json.Marshal(tt.args)

			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(b))
			router.ServeHTTP(w, req)

			assert.NotEqual(t, 201, w.Code, tt.name)
      exp := `{"success": false, "message": "Invalid request; ", "data": null}`

			if !reflect.DeepEqual(w.Body.String(), exp) {
				t.Errorf("RegisterUser() = %s , want=%s", w.Body.String(), exp)
			}
		})
	}
}