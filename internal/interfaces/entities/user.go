package entities

import (
	"app/internal/models"
	"time"
)

type UserData struct {
	Fullname    string    `json:"fullname"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Birthday    time.Time `json:"birthday"`
	LatestLogin time.Time `json:"latest_login"`
}

func (u UserData) ToModel() (*models.User) {
  return &models.User{
    Fullname: u.Fullname,
    PhoneNumber: u.PhoneNumber,
    Email: u.Email,
    Username: u.Username,
    Password: u.Password,
    Birthday: u.Birthday,
  }
}
