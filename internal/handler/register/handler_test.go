package register

import (
	"app/internal/interfaces/entities"
	"app/internal/interfaces/usecases/user/mocks"
	"app/internal/models"
	"app/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	sampleUserData entities.UserData = entities.UserData{
		Fullname:    "Nguyen Quoc Dat",
		PhoneNumber: "123456790",
		Email:       "datshiro@gmail.com",
		Username:    "datshiro",
		Password:    "sTr0ngP@ssword",
	}

	sampleUser *models.User = &models.User{
		Fullname:    "Nguyen Quoc Dat",
		PhoneNumber: "123456790",
		Email:       "datshiro@gmail.com",
		Username:    "datshiro",
		Password:    "sTr0ngP@ssword",
	}
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	mock := new(mocks.UserUsecase)
	mock.
		On("RegisterUser", context.Background(), sampleUserData).
		Return(sampleUser, nil)
	r.POST("/register", NewHandler(mock).Handle)
	return r
}

func TestRegisterRoute(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name     string
		args     entities.UserData
		expected *models.User
		wantErr  bool
	}{
		{
			"Register Success",
			entities.UserData{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "123456790",
				Email:       "datshiro@gmail.com",
				Username:    "datshiro",
				Password:    "sTr0ngP@ssword",
			},
			&models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "123456790",
				Email:       "datshiro@gmail.com",
				Username:    "datshiro",
				Password:    "sTr0ngP@ssword",
			},
			false,
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		b, err := json.Marshal(tt.args)
		assert.NoError(t, err, "Marshal request body failed")

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(b))
		router.ServeHTTP(w, req)

		expectedResponse := utils.ResponseObject{Success: true, Data: tt.expected}
		expected, err := json.Marshal(expectedResponse)
		assert.NoError(t, err, "Marshal request body failed")

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, bytes.NewBuffer(expected).String(), w.Body.String())
	}
}
