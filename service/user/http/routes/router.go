package routes

import (
	"github.com/GoCloudstorage/GoCloudstorage/service/user/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func RouterInit() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())

	user := app.Group("/user")
	{
		user.Post("/create", http.UserRegister)
		user.Post("/login", http.UserLogin)
	}
	return app
}
