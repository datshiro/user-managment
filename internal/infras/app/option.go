package app

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
)

type Opts struct {
	Port                  int    `env:"PORT" envDefault:"3333"`
	ApiPrefix             string `env:"API_PREFIX" envPrefix:"api"`
	DatabaseEnable        bool   `env:"DATABASE_ENABLE" envDefault:"false"`
	RedisEnable           bool   `env:"REDIS_ENABLE" envDefault:"false"`
	RedisConnectionString string `env:"REDIS_CONNECTION_STRING"`
}

func defaultOpts() Opts {
	c := Opts{
		Port:                  3000,
		RedisConnectionString: "redis://0.0.0.0:6379",
		ApiPrefix:             "api",
		DatabaseEnable:        false,
		RedisEnable:           false,
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

func WithDatabase(isConnectDatabase bool) OptFunc {
	return func(c *Opts) {
		c.DatabaseEnable = isConnectDatabase
	}
}
