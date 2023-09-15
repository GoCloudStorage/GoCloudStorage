package handler

import (
	"context"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/user/user"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/xrpc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type API struct {
	UserRPC user.UserServiceClient
}

func (a *API) InitGrpc() {
	// add storage rpc client
	client, err := xrpc.GetGrpcClient(xrpc.Config{
		Domain:          "user",
		Endpoints:       []string{"localhost:50001"},
		BackoffInterval: 0,
		MaxAttempts:     0,
	},
		user.NewUserServiceClient,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	a.UserRPC = client.NewSession()
}

func (a *API) RouterInit() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())

	user := app.Group("/user")
	{
		user.Post("/create", a.UserRegister)
		user.Post("/login", a.UserLogin)
	}
	return app
}

var api API

func InitAPI(ctx context.Context) {
	var (
		addr = fmt.Sprintf("%s:%s", opt.Cfg.CloudStorage.Host, opt.Cfg.CloudStorage.Port)
		app  = api.RouterInit()
	)
	api.InitGrpc()
	go func() {
		logrus.Infof("Start fiber webserver, addr: %s", addr)
		if err := app.Listen(addr); err != nil {
			logrus.Panicf("%s listen address %v fail, err: %v", opt.Cfg.CloudStorage.Name, addr, err)
		}
	}()
	select {
	case <-ctx.Done():
		_ = app.Shutdown()
	}
}
