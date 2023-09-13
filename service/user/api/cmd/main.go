package main

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"github.com/GoCloudstorage/GoCloudstorage/service/user/http/routes"
	"github.com/GoCloudstorage/GoCloudstorage/service/user/rpc/server"

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
	server.ClientInit()
	defer server.Conn.Close()
	app := routes.RouterInit()
	app.Listen(":8080")
	<-ctx.Done()
	logrus.Warnf("cloud user service stop by ctx in 3s...")
	<-time.After(time.Second * 3)
}
