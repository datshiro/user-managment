package models

import (
	"app/internal/consts"
	"app/internal/utils"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
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

// Hash password before insert into database
func (u *User) BeforeCreate(tx *gorm.DB) (err error)  {
  u.Password, err = utils.HashPassword(u.Password)

  if err != nil {
    if errors.Is(err, bcrypt.ErrPasswordTooLong) {
      err = consts.NewError(fmt.Errorf("password too long; %v", err))
      return
    }
    if errors.Is(err, bcrypt.ErrHashTooShort) {
      err = consts.NewError(fmt.Errorf("password too short; %v", err))
      return
    }
    err = consts.NewError(fmt.Errorf("failed to hash password; %v", err))
  }
  return
}
