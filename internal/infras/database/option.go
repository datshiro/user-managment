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
	dialector gorm.Dialector
}
type OptFunc func(*Opts)

func loadOpts() Opts {
	c := Opts{}

	if err := env.Parse(&c); err != nil {
		log.Fatalf("%+v\n", err)
	}
  fmt.Printf("%+v\n", c)
	return c

}

func (opt *Opts) MakeConnect() (*gorm.DB, error) {
	db, err := gorm.Open(opt.dialector, &gorm.Config{})
	return db, err
}
