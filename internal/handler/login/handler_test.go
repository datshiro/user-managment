package login

import (
	"app/internal/interfaces/app/middlewares"
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
	r.POST("/", NewHandler(uc).Handle)
	return r
}

func Test_loginHandler_Handle(t *testing.T) {
	type fields struct {
		OnMethod         string // Mock usecase method name
		Arguments        []interface{}
		ReturnValueUser  *models.User
		ReturnValueError error
	}
	type args struct {
		request
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		expected    interface{}
		expectedErr error
	}{
		{
			name: "success login with email",
			fields: fields{
				OnMethod:         "LoginWithEmail",
				Arguments:        []interface{}{context.Background(), "datshiro@gmail.com", "sTr0ngP@ssword"},
				ReturnValueUser:  &models.User{Email: "datshiro@gmail.com", Password: "sTr0ngP@ssword"},
				ReturnValueError: nil,
			},
			args: args{
				request: request{
					Email:    "datshiro@gmail.com",
					Password: "sTr0ngP@ssword",
				},
			},
			expected: response{
				Email: "datshiro@gmail.com",
			},
			expectedErr: nil,
		},
		{
			name: "success login with phone",
			fields: fields{
				OnMethod:         "LoginWithPhone",
				Arguments:        []interface{}{context.Background(), "0123456789", "sTr0ngP@ssword"},
				ReturnValueUser:  &models.User{PhoneNumber: "0123456789", Password: "sTr0ngP@ssword"},
				ReturnValueError: nil,
			},
			args: args{
				request: request{
					PhoneNumber: "0123456789",
					Password:    "sTr0ngP@ssword",
				},
			},
			expected: response{
				PhoneNumber: "0123456789",
			},
			expectedErr: nil,
		},
		{
			name: "success login with username",
			fields: fields{
				OnMethod:         "LoginWithUsername",
				Arguments:        []interface{}{context.Background(), "datshiro@gmail.com", "sTr0ngP@ssword"},
				ReturnValueUser:  &models.User{Username: "datshiro@gmail.com", Password: "sTr0ngP@ssword"},
				ReturnValueError: nil,
			},
			args: args{
				request: request{
					Username: "datshiro@gmail.com",
					Password: "sTr0ngP@ssword",
				},
			},
			expected: response{
				Username: "datshiro@gmail.com",
			},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := new(mocks.UserUsecase)
			mock.On(tt.fields.OnMethod, tt.fields.Arguments...).Return(tt.fields.ReturnValueUser, tt.fields.ReturnValueError)
			router := setupRouter(mock)

			w := httptest.NewRecorder()
			b, err := json.Marshal(tt.args.request)
			assert.NoError(t, err, "Marshal request body failed")

			req, err := http.NewRequest("POST", "/", bytes.NewBuffer(b))
			assert.NoError(t, err, "Create test request failed")
			router.ServeHTTP(w, req)

			if tt.expectedErr == nil {
				// Parse success result
				expectedResponse := utils.ResponseObject{Success: true, Data: tt.expected}
				expected, err := json.Marshal(expectedResponse)
				assert.NoError(t, err, "Marshal request body failed")

				assert.Equal(t, 200, w.Code, "HTTP Status Code not 200")
				assert.Equal(t, bytes.NewBuffer(expected).String(), w.Body.String(), "Response data not as expected")

			} else {
				assert.NotEqual(t, 200, w.Code, "HTTP Status Code is 200")
				assert.Equal(t, fmt.Sprintf("%q", tt.expectedErr.Error()), w.Body.String(), "Error response not as expected")
			}

		})
	}
}
