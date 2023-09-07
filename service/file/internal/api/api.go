package api

import (
	"context"
	"fmt"
	"github.com/GoCloudstorage/GoCloudstorage/opt"
	"github.com/GoCloudstorage/GoCloudstorage/service/file/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func registerAPI() *fiber.App {
	app := fiber.New()
	api := app.Group("/file")
	{
		api.Post("/", handler.PreUpload)
		api.Get("/", handler.GetAll)
		api.Get("/:id", handler.PreDownload)
	}
	return app
}

func InitAPI(ctx context.Context) {
	var (
		addr = fmt.Sprintf("%s:%s", opt.Cfg.CloudStorage.Host, opt.Cfg.CloudStorage.Port)
		app  = registerAPI()
	)
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
