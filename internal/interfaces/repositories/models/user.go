package models

import (
	"app/internal/consts"
	"app/internal/utils"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
  ErrPasswordTooLong = errors.New("password: too long for hashing")
  ErrPasswordTooShort = errors.New("password: too short for hashing")
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
	u.Password, err = utils.HashPassword(u.Password)

	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			err = consts.NewError(ErrPasswordTooLong)
			return
		}
		if errors.Is(err, bcrypt.ErrHashTooShort) {
			err = consts.NewError(ErrPasswordTooShort)
			return
		}
		err = consts.NewError(ErrPasswordHash).WithRootCause(err)
	}
	return
}

func (u *User) ComparePassword(password string) bool {
  return utils.CheckPasswordHash(password, u.Password)
}
