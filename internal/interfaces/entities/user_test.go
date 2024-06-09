package entities

import (
	"app/internal/models"
	"reflect"
	"testing"
	"time"
)

func TestUserData_ToModel(t *testing.T) {
	now := time.Now()
	type fields struct {
		Fullname    string
		PhoneNumber string
		Email       string
		Username    string
		Password    string
		Birthday    time.Time
		LatestLogin time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   *models.User
	}{
		{
			name: "Success To Model",
			fields: fields{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "123456789",
				Email:       "datshiro@gmail.com",
				Username:    "datshiro",
				Password:    "my_P$ssword",
				Birthday:    time.Date(1996, time.July, 19, 0, 0, 0, 0, time.UTC),
				LatestLogin: now,
			},
			want: &models.User{
				Fullname:    "Nguyen Quoc Dat",
				PhoneNumber: "123456789",
				Email:       "datshiro@gmail.com",
				Username:    "datshiro",
				Password:    "my_P$ssword",
				Birthday:    time.Date(1996, time.July, 19, 0, 0, 0, 0, time.UTC),
				LatestLogin: now,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := UserData{
				Fullname:    tt.fields.Fullname,
				PhoneNumber: tt.fields.PhoneNumber,
				Email:       tt.fields.Email,
				Username:    tt.fields.Username,
				Password:    tt.fields.Password,
				Birthday:    tt.fields.Birthday,
				LatestLogin: tt.fields.LatestLogin,
			}
			if got := u.ToModel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserData.ToModel() = %v, want %v", got, tt.want)
			}
		})
	}
}
