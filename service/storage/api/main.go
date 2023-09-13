package main

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/storage_engine"
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
	model.Init()
	storage_engine.InitClient(opt.Cfg.Storage.Type, opt.Cfg.Storage.Endpoint, opt.Cfg.Storage.AccessKeyID, opt.Cfg.Storage.SecretAccessKey, opt.Cfg.Storage.BucketName, opt.Cfg.Storage.UseSSL)
	api.InitAPI(ctx)
	<-ctx.Done()
	logrus.Warnf("cloud file service stop by ctx in 3s...")
	<-time.After(time.Second * 3)
}
