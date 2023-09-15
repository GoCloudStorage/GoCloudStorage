package main

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
<<<<<<< HEAD
	"github.com/GoCloudstorage/GoCloudstorage/pkg/storage_engine"
=======
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/redis"
>>>>>>> f5f05860dc07a675e4e61571dfb88bb9103fede2
	"github.com/GoCloudstorage/GoCloudstorage/service/storage/api/internal/handler"
	"github.com/GoCloudstorage/GoCloudstorage/service/storage/model"
	"github.com/sirupsen/logrus"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	opt.InitConfig()
	pg.Init(opt.Cfg.Pg.Host, opt.Cfg.Pg.User, opt.Cfg.Pg.Password, opt.Cfg.Pg.DBName, opt.Cfg.Pg.Port)
	redis.Init(opt.Cfg.Redis.Addr, opt.Cfg.Redis.Password, opt.Cfg.Redis.DB)
	model.Init()
<<<<<<< HEAD
	storage_engine.InitClient(opt.Cfg.Storage.Type, opt.Cfg.Storage.Endpoint, opt.Cfg.Storage.AccessKeyID, opt.Cfg.Storage.SecretAccessKey, opt.Cfg.Storage.BucketName, opt.Cfg.Storage.UseSSL)
	api.InitAPI(ctx)
=======
	SvcConfig := opt.Cfg.StorageService
	handler.InitAPI(ctx, SvcConfig.Name, SvcConfig.Host, SvcConfig.Port)
>>>>>>> f5f05860dc07a675e4e61571dfb88bb9103fede2
	<-ctx.Done()
	logrus.Warnf("cloud file service stop by ctx in 3s...")
	<-time.After(time.Second * 3)
}
