package bootstrap

import (
	"go-chat-app/app/websocket"
	"go-chat-app/pkg/database"
	"go-chat-app/pkg/env"
	"go-chat-app/pkg/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func NewApplication() *fiber.App {
	env.SetupEnvFile()
	database.SetupDatabase()
	database.SetupMongoDb()
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Get("/dashboard", monitor.New())
	go websocket.ServeWsMessage(app)
	router.InstallRouter(app)
	return app
}
