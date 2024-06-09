package app

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
)

type Opts struct {
	Port              int    `env:"PORT" envDefault:"3333" json:"port"`
	RedisUrl          string `env:"REDIS_URL" json:"redis_url"`
	ApiPrefix         string `env:"API_PREFIX" envPrefix:"api" json:"api_prefix"`
	IsConnectDatabase bool   `env:"IS_CONNECT_DATABASE" envDefault:"false"`
}

func defaultOpts() Opts {
	c := Opts{
		Port:      3000,
		RedisUrl:  "redis://0.0.0.0:6379",
		ApiPrefix: "api",
    IsConnectDatabase: false,
	}

	if err := env.Parse(&c); err != nil {
		log.Fatalf("%+v\n", err)
	}
	fmt.Printf("%+v\n", c)

	return c
}

type OptFunc func(*Opts)

func WithPort(port int) OptFunc {
	return func(c *Opts) {
		c.Port = port
	}
}
