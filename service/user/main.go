package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"os/signal"
	"syscall"
	"time"
	"work-space/opt"
	"work-space/tools/db/pg"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	opt.InitConfig()
	pg.Init(opt.Cfg.Pg.Host, opt.Cfg.Pg.DBName, opt.Cfg.Pg.User, opt.Cfg.Pg.Password, opt.Cfg.Pg.Port)
	api.InitAPI(ctx)
	<-ctx.Done()
	logrus.Warnf("cloud storage service stop by ctx in 3s...")
	<-time.After(time.Second * 3)
}
