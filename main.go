package main

import (
	"errors"
	"fmt"
	"github.com/freegle/iznik-server-go/database"
	"github.com/freegle/iznik-server-go/router"
	"github.com/freegle/iznik-server-go/user"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 8)

	// This runs on the server where the timezone should be set to UTC.  Make sure that's also true when we're
	// running in development.
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc

	app := fiber.New(fiber.Config{
		ReadBufferSize:  8192,
		WriteBufferSize: 8192,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Map this to a standardised error response.
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			return ctx.Status(code).JSON(fiber.Map{
				"error":   code,
				"message": err.Error(),
			})
		},
	})

	// Use compression unless we're inside the Docker environment.
	if strings.Index(".localhost", os.Getenv("USER_SITE")) < 0 {
		app.Use(compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}))
	}

	// Enable CORS - we don't care who uses the API.  Set MaxAge so that OPTIONS preflight requests are cached, which
	// reduces the number of them and hence increases performance.
	app.Use(cors.New(cors.Config{
		MaxAge: 86400,
	}))

	database.InitDatabase()

	app.Use(database.NewPingMiddleware(database.Config{}))

	router.SetupRoutes(app)

	// Add our middleware to check for a valid JWT. Do this after the ping middleware - I think the middleware
	// execution order is in the order that they're added.
	app.Use(user.NewAuthMiddleware(user.Config{}))

	// We can signal to stop using SIGINT.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	serverShutdown := make(chan struct{})

	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
		serverShutdown <- struct{}{}
	}()

	app.Listen(":8192")

	<-serverShutdown

	fmt.Println("...exiting")
}
