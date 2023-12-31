package redis

import (
	"context"

	"errors"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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

// SetLock 给redis添加分布式锁
func SetLock(ctx context.Context, lockKey string, lockTimeout time.Duration) bool {
	ok, err := Client.SetNX(ctx, lockKey, opt.Cfg.Redis.UniqueValue, lockTimeout).Result()
	if err != nil {
		logrus.Error("redis set nx err:", err)
		return false
	}
	return ok
}

// ReleaseLock 解分布式锁
func ReleaseLock(ctx context.Context, lockKey string) {
	// 使用 Lua 脚本来释放锁
	luaScript := `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
`
	_, err := Client.Eval(ctx, luaScript, []string{lockKey}, opt.Cfg.Redis.UniqueValue).Result()
	if err != nil {
		logrus.Error("redis releaseLock err:", err)
	}
	return
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
