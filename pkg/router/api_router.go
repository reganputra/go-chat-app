package router

import (
	"go-chat-app/app/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type ApiRouter struct {
}

func (a ApiRouter) InstallRouter(app *fiber.App) {
	api := app.Group("/api", limiter.New())
	api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Hello from api",
		})
	})
	userGroup := api.Group("/user")
	userV1 := userGroup.Group("/v1")
	userV1.Post("/register", controllers.RegisterUser)
	userV1.Post("/login", controllers.LoginUser)
}
func NewApiRouter() *ApiRouter {
	return &ApiRouter{}
}
