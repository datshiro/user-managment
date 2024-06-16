package redis

import "github.com/redis/go-redis/v9"

func NewRedis(opts ...OptFunc) *redis.Client {
	o := Opts{}
	connOpts, err := redis.ParseURL(o.ConnectionString)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(connOpts)
}
