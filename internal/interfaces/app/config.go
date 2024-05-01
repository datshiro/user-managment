package app

type Config struct {
	Port     int    `env:"PORT"`
	DbUrl    string `env:"DB_URL"`
	RedisUrl string `env:"REDIS_URL"`
}
