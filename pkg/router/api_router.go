package router

import (
	"go-chat-app/app/controllers"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type ApiRouter struct {
}

func (a ApiRouter) InstallRouter(app *fiber.App) {
	api := app.Group("/api", limiter.New(limiter.Config{
		Max:        50,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(ctx *fiber.Ctx) string {
			return ctx.IP()
		},
	}))
	api.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Hello from api",
		})
	})
	userGroup := api.Group("/user")
	userV1 := userGroup.Group("/v1")
	userV1.Post("/register", controllers.RegisterUser)
	userV1.Post("/login", controllers.LoginUser)
	userV1.Delete("/logout", AuthMiddleware, controllers.LogoutUser)
	userV1.Put("/refresh-token", MiddlewareRefreshToken, controllers.RefreshToken)

	messageGroup := api.Group("/message")
	messageV1 := messageGroup.Group("/v1")
	messageV1.Get("/history", AuthMiddleware, controllers.GetMessagesHistory)
}
func NewApiRouter() *ApiRouter {
	return &ApiRouter{}
}
