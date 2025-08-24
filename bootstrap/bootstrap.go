package bootstrap

import (
	"go-chat-app/app/websocket"
	"go-chat-app/pkg/database"
	"go-chat-app/pkg/env"
	"go-chat-app/pkg/router"
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
	"go.elastic.co/apm"
)

func NewApplication() *fiber.App {
	env.SetupEnvFile()
	SetupLogFile()

	database.SetupDatabase()
	database.SetupMongoDb()

	apm.DefaultTracer.Service.Name = "go-chat-app"
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Get("/dashboard", monitor.New())

	go websocket.ServeWsMessage(app)

	router.InstallRouter(app)
	return app
}

func SetupLogFile() {
	logFile, err := os.OpenFile("./logs/chat_message.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
