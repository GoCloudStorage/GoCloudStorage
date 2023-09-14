package redis

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

var Client *redis.Client

func Init(addr string, password string, db int) {
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

// 给redis添加分布式锁
func SetLock(ctx context.Context, lockKey string, lockTimeout time.Duration) bool {
	ok, err := Client.SetNX(ctx, lockKey, opt.Cfg.Redis.UniqueValue, lockTimeout).Result()
	if err != nil {
		logrus.Error("redis set nx err:", err)
		return false
	}
	return ok
}

// 解锁
func ReleaseLock(ctx context.Context, lockKey string) bool {
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
		return false
	}
	return true
}
