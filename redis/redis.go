package redis

import (
	"github.com/go-redis/redis/v7"
)

type RedisConf struct {
	Host string
	Pwd  string
	Db   int
}

func NewConn(r *RedisConf) *redis.Client { // in huge app it will be global config object
	return redis.NewClient(&redis.Options{
		Addr:     r.Host,
		Password: r.Pwd, // no password set
		DB:       r.Db,  // use default DB
	})
}
