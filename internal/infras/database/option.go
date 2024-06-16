package database

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"gorm.io/gorm"
)

type Opts struct {
	DB        string `env:"DATABASE_DB"`
	User      string `env:"DATABASE_USER"`
	Password  string `env:"DATABASE_PASSWORD"`
	Host      string `env:"DATABASE_HOST"`
	Port      string `env:"DATABASE_PORT"`
	SSLMode   string `env:"DATABASE_SSL_MODE" default:"disable"`
	TimeZone  string `env:"DATABASE_TIME_ZONE" default:"Asia/Shanghai"`
	Dialector gorm.Dialector
}
type OptFunc func(*Opts)

func LoadOpts() Opts {
	c := Opts{}

	if err := env.Parse(&c); err != nil {
		log.Fatalf("%+v\n", err)
	}
	fmt.Printf("%+v\n", c)
	return c

}

// db_type must be postgres or other database type
func (opt *Opts) GetDSN() string {
	dsn := "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s"
	return fmt.Sprintf(dsn, opt.Host, opt.User, opt.Password, opt.DB, opt.Port, opt.SSLMode, opt.TimeZone)
}

func WithUser(user string) OptFunc {
  return func(o *Opts) {
    o.User = user
  }
}

func WithPassword(pass string) OptFunc {
  return func(o *Opts) {
    o.Password = pass
  }
}

func WithDB(db_name string) OptFunc {
  return func(o *Opts) {
    o.DB = db_name
  }
}
