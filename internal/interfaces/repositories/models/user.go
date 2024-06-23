package models

import (
	"app/internal/utils"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
  ErrPasswordTooLong = errors.New("password: too long")
  ErrPasswordTooShort = errors.New("password: too short")
  ErrPasswordHash = errors.New("password: hash failed")
)

type User struct {
	gorm.Model
	Fullname    string
	PhoneNumber string `gorm:"not_null"`
	Email       string `gorm:"unique"`
	Password    string `gorm:"not_null"`
	Birthday    time.Time
	LatestLogin time.Time
}

// Hash password before insert into database
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
  if len(u.Password) < 7 {
    return ErrPasswordTooShort
  }
	u.Password, err = utils.HashPassword(u.Password)

	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return ErrPasswordTooLong
		}
	  return err
	}
	return nil
}

func (u *User) ComparePassword(password string) bool {
  return utils.CheckPasswordHash(password, u.Password)
}
