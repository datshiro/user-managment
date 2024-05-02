package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
  Fullname    string 
  PhoneNumber string `gorm:"unique"`
  Email       string `gorm:"unique"`
  Username    string `gorm:"unique"`
	Password    string
	Birthday    time.Time
	LatestLogin time.Time
}
