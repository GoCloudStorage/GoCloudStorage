package main

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"github.com/GoCloudstorage/GoCloudstorage/service/file/api/internal/handler"
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
<<<<<<< HEAD

	handler.InitAPI(ctx)
=======
	model.Init()
	SvcConfig := opt.Cfg.FileService
	handler.InitAPI(ctx, SvcConfig.Name, SvcConfig.Host, SvcConfig.Port)
>>>>>>> f5f05860dc07a675e4e61571dfb88bb9103fede2
	<-ctx.Done()
	logrus.Warnf("cloud file service stop by ctx in 3s...")
	<-time.After(time.Second * 3)
}
