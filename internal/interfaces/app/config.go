package app

import (
	"app/internal/infras/database"
)

type Config struct {
	Port      int               `env:"PORT"`
	RedisUrl  string            `env:"REDIS_URL"`
	ApiPrefix string            `env:"API_PREFIX"`
	DbConfig  database.DBConfig `envPrefix:"DB_"`
}
