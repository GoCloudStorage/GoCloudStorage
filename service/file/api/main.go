package main

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"github.com/GoCloudstorage/GoCloudstorage/service/file/api/internal/handler"
	"github.com/GoCloudstorage/GoCloudstorage/service/file/model"
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
	SvcConfig := opt.Cfg.FileService
	handler.InitAPI(ctx, SvcConfig.Name, SvcConfig.Host, SvcConfig.Port)
	<-ctx.Done()
	logrus.Warnf("cloud file logic stop by ctx in 3s...")
	<-time.After(time.Second * 3)
}
