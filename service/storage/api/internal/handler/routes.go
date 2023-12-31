package handler

import (
	"context"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/xrpc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

type API struct {
	storageRPCClient storage.StorageClient
}

func (a *API) InitGrpc() {
	// add local rpc client
	client, err := xrpc.GetGrpcClient(
		xrpc.Config{
			Domain:          opt.Cfg.StorageRPC.Domain,
			Endpoints:       opt.Cfg.StorageRPC.Endpoints,
			BackoffInterval: 0,
			MaxAttempts:     0,
		},
		storage.NewStorageClient,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		panic(err)
	}
	a.storageRPCClient = client.NewSession()
}

func (a *API) registerAPI() *fiber.App {
	app := fiber.New()
	api := app.Group("")
	api.Use(cors.New())
	{
		api.Get("/", func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(http.StatusOK)
		})
		//api.Put("/upload/:token", a.upload)
		app.Put("/upload/:token", a.upload)
		api.Get("/download/:token", a.download)

	}
	return app
}

var api API

func InitAPI(ctx context.Context, name string, host string, port string) {
	var (
		addr = fmt.Sprintf("%s:%s", host, port)
		app  = api.registerAPI()
	)
	api.InitGrpc()
	go func() {
		logrus.Infof("Start fiber webserver, addr: %s", addr)
		if err := app.Listen(addr); err != nil {
			logrus.Panicf("%s listen address %v fail, err: %v", name, addr, err)
		}
	}()
	select {
	case <-ctx.Done():
		_ = app.Shutdown()
	}
}
