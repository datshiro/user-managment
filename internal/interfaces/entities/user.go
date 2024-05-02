package entities

import "time"

type UserData struct {
	Fullname    string    `json:"fullname"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Birthday    time.Time `json:"birthday"`
	LatestLogin time.Time `json:"latest_login"`
}
