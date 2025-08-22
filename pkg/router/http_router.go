package router

import (
	"go-chat-app/app/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
)

type HttpRouter struct {
}

func (h HttpRouter) InstallRouter(app *fiber.App) {
	group := app.Group("", cors.New(), csrf.New())
	group.Get("/", controllers.RenderAuth)
	group.Get("/auth", controllers.RenderAuth)
	group.Get("/chat", controllers.RenderChat)
	group.Get("/dashboard-ui", controllers.RenderUI)
}

func NewHttpRouter() *HttpRouter {
	return &HttpRouter{}
}
