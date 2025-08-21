package router

import (
	"go-chat-app/app/repositories"
	"go-chat-app/pkg/jwt"
	"go-chat-app/pkg/response"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		log.Println("No auth token")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	// Extract token from "Bearer <token>" format
	var auth string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		auth = authHeader[7:]
	} else {
		auth = authHeader
	}

	_, err := repositories.GetUserSession(ctx.Context(), auth)
	if err != nil {
		log.Println("failed to get user session", err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	claims, err := jwt.ValidateToken(ctx.Context(), auth)
	if err != nil {
		log.Println("Invalid token")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		log.Println("Token expired", claims.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	ctx.Set("username", claims.Username)
	ctx.Set("full_name", claims.FullName)
	return ctx.Next()
}
