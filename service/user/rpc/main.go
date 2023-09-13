package main

import (
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/user/user"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"github.com/GoCloudstorage/GoCloudstorage/service/user/rpc/internal/db/model"
	"github.com/GoCloudstorage/GoCloudstorage/service/user/rpc/internal/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func main() {
	opt.InitConfig()
	pg.Init(opt.Cfg.Pg.Host, opt.Cfg.Pg.User, opt.Cfg.Pg.Password, opt.Cfg.Pg.DBName, opt.Cfg.Pg.Port)
	model.Init()
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logrus.Error("start user rpc listener err:", err)
		return
	}
	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, service.GetUserSrv())
	err = s.Serve(listener)
	if err != nil {
		logrus.Error("start user rpc err:", err)
		return
	}
}
