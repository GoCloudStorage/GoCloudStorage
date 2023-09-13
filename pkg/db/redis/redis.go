package redis

import (
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init(addr string, password string, db int) {
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}
