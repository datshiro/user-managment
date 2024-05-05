package register

import (
	"app/internal/consts"
	"app/internal/interfaces/app/middlewares"
	"app/internal/interfaces/entities"
	"app/internal/interfaces/usecases/user"
	mocks "app/internal/mocks/usecases/user"
	"app/internal/models"
	"app/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(uc user.UserUsecase) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middlewares.ErrorHandlerMiddleware())
	r.POST("/register", NewHandler(uc).Handle)
	return r
}

func TestRegisterRoute(t *testing.T) {
	tests := []struct {
		name        string
		args        entities.UserData
		expected    *models.User
		expectedErr error
	}{
		{
			"success with phone only",
			entities.UserData{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "123456790",
				Email:       "",
				Username:    "",
				Password:    "sTr0ngP@ssword",
			},
			&models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "123456790",
				Email:       "",
				Username:    "",
				Password:    "sTr0ngP@ssword",
			},
			nil,
		},
		{
			"success: with username only ",
			entities.UserData{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "",
				Email:       "",
				Username:    "datshiro",
				Password:    "sTr0ngP@ssword",
			},
			&models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "",
				Email:       "",
				Username:    "datshiro",
				Password:    "sTr0ngP@ssword",
			},
			nil,
		},
		{
			"success: with email only ",
			entities.UserData{
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
			nil,
		},
		{
			"failed: missing account info",
			entities.UserData{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "",
				Email:       "",
				Username:    "",
				Password:    "sTr0ngP@ssword",
			},
			&models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "123456790",
				Email:       "datshiro@gmail.com",
				Username:    "datshiro",
				Password:    "sTr0ngP@ssword",
			},
			consts.ErrMissingCredentialInfo,
		},
		{
			"failed: missing password",
			entities.UserData{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "123456790",
				Email:       "datshiro@gmail.com",
				Username:    "datshiro",
				Password:    "",
			},
			&models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "123456790",
				Email:       "datshiro@gmail.com",
				Username:    "datshiro",
				Password:    "",
			},
			consts.ErrPasswordCannotBeEmpty,
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

			if tt.expectedErr == nil {
				// Parse success result
				expectedResponse := utils.ResponseObject{Success: true, Data: tt.expected}
				expected, err := json.Marshal(expectedResponse)
				assert.NoError(t, err, "Marshal request body failed")

				assert.Equal(t, 200, w.Code, tt.name)
				assert.Equal(t, bytes.NewBuffer(expected).String(), w.Body.String(), tt.name)
			} else {
				assert.NotEqual(t, 200, w.Code, tt.name)
				assert.Equal(t, fmt.Sprintf("%q", tt.expectedErr.Error()), w.Body.String(), tt.name)
			}
		})
	}
}
