package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

var Client *redis.Client

var (
	Nil = errors.New("NotRecord")
)

func Init(addr string, password string, db int) {
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func Get(ctx context.Context, key string) (string, error) {
	cmd := Client.Get(ctx, key)
	if cmd.Err() == redis.Nil {
		return "", Nil
	} else if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Result()
}

func SetEx(ctx context.Context, key string, value any, expire time.Duration) error {
	return Client.SetEx(ctx, key, value, expire).Err()
}
