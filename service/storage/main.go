package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"os/signal"
	"syscall"
	"time"
	"work-space/opt"
	"work-space/service/storage/internal/api"
	"work-space/service/storage/internal/model"
	"work-space/tools/db/pg"
	"work-space/tools/storage/minio"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	opt.InitConfig()
	pg.Init(opt.Cfg.Pg.Host, opt.Cfg.Pg.User, opt.Cfg.Pg.Password, opt.Cfg.Pg.DBName, opt.Cfg.Pg.Port)
	minio.Init(opt.Cfg.Storage.Endpoint, opt.Cfg.Storage.AccessKeyID, opt.Cfg.Storage.SecretAccessKey, opt.Cfg.Storage.UseSSL)
	model.Migrator()
	api.InitAPI(ctx)
	<-ctx.Done()
	logrus.Warnf("cloud storage service stop by ctx in 3s...")
	<-time.After(time.Second * 3)
}
