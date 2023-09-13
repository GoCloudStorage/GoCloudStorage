package main

import (
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/file"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/pg"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/db/redis"
	"github.com/GoCloudstorage/GoCloudstorage/service/file/model"
	"github.com/GoCloudstorage/GoCloudstorage/service/file/rpc/internal/logic"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func main() {
	opt.InitConfig()
	pg.Init(opt.Cfg.Pg.Host, opt.Cfg.Pg.User, opt.Cfg.Pg.Password, opt.Cfg.Pg.DBName, opt.Cfg.Pg.Port)
	redis.Init(opt.Cfg.Redis.Addr, opt.Cfg.Redis.Password, opt.Cfg.Redis.DB)
	model.Init()

	listener, err := net.Listen("tcp", opt.Cfg.FileRPC.Endpoints[0])
	if err != nil {
		logrus.Error("start storage rpc listener err:", err)
		return
	}
	s := grpc.NewServer()
	file.RegisterFileServer(s, &logic.FileServer{})
	logrus.Infof("start listen %s", opt.Cfg.FileRPC.Endpoints[0])
	err = s.Serve(listener)
	if err != nil {
		logrus.Error("start storage rpc err:", err)
		return
	}
}
