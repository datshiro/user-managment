package redis

type Opts struct {
	ConnectionString        string `env:"REDIS_CONNECTION_STRING"`
}
type OptFunc func(*Opts)
