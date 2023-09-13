package main

import (
	"context"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/user/user"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"github.com/GoCloudstorage/GoCloudstorage/service/user/http/routes"
	"github.com/GoCloudstorage/GoCloudstorage/service/user/internal/db/model"
	"github.com/GoCloudstorage/GoCloudstorage/service/user/internal/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
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
	routes.InitAPI(ctx)

	server := grpc.NewServer()
	defer server.Stop()
	user.RegisterUserServiceServer(server, service.GetUserSrv())
	lis, err := net.Listen("tcp", ":50051") //addr需要通过etcd拿到具体addr
	if err != nil {
		panic(err)
	}

	if err = server.Serve(lis); err != nil {
		panic(err)
	}
	<-ctx.Done()
	logrus.Warnf("cloud file service stop by ctx in 3s...")
	<-time.After(time.Second * 3)
}
