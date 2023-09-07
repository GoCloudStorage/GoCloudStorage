package api

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"work-space/opt"
	"work-space/service/storage/internal/handler"
)

func registerAPI() *fiber.App {
	app := fiber.New()
	api := app.Group("/storage")
	{
		api.Post("/", handler.Upload)
		api.Get("/", handler.GetAll)
		api.Get("/:id", handler.Download)

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
