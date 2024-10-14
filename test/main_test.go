package test

import (
	"github.com/freegle/booktastic-server-go/database"
	"github.com/freegle/booktastic-server-go/router"
	"github.com/freegle/booktastic-server-go/user"
	"github.com/gofiber/fiber/v2"
)

var app *fiber.App

func init() {
	app = fiber.New()
	app.Use(user.NewAuthMiddleware(user.Config{}))
	database.InitDatabase()
	router.SetupRoutes(app)
}

func getApp() *fiber.App {
	// We use this so that we only initialise fiber once.
	return app
}
