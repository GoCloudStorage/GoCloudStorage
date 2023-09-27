package main

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"github.com/GoCloudstorage/GoCloudstorage/service/user/http/handler"
	"github.com/sirupsen/logrus"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logrus.SetReportCaller(true)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	opt.InitConfig()
	pg.Init(opt.Cfg.Pg.Host, opt.Cfg.Pg.User, opt.Cfg.Pg.Password, opt.Cfg.Pg.DBName, opt.Cfg.Pg.Port)

	userConfig := opt.Cfg.UserService

	handler.InitAPI(ctx, userConfig.Name, userConfig.Host, userConfig.Port)

	<-ctx.Done()
	logrus.Warnf("cloud user logic stop by ctx in 3s...")
	<-time.After(time.Second * 3)
}
