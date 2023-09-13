package handler

import (
	"context"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/pb/storage"
	"github.com/GoCloudstorage/GoCloudstorage/pkg/xrpc"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type API struct {
	storageRPC storage.StorageClient
}

func (a *API) InitGrpc() {
	// add storage rpc client
	client, err := xrpc.GetGrpcClient(xrpc.Config{
		Domain:          "localhost:8001",
		Endpoints:       nil,
		BackoffInterval: 0,
		MaxAttempts:     0,
	}, storage.NewStorageClient)
	if err != nil {
		panic(err)
	}
	a.storageRPC = client.NewSession()

}

func (a *API) registerAPI() *fiber.App {
	app := fiber.New()

	api := app.Group("/file")
	{
		api.Post("/", preUpload)
		//api.Get("/", GetAll)
		api.Get("/:id", a.preDownload)
	}
	return app
}

var api API

func InitAPI(ctx context.Context) {
	var (
		addr = fmt.Sprintf("%s:%s", opt.Cfg.CloudStorage.Host, opt.Cfg.CloudStorage.Port)
		app  = api.registerAPI()
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
