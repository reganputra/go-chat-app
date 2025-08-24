package router

import (
	"go-chat-app/app/repositories"
	"go-chat-app/pkg/jwt"
	"go-chat-app/pkg/response"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
)

func AuthMiddleware(ctx *fiber.Ctx) error {

	span, spanCtx := apm.StartSpan(ctx.Context(), "AuthMiddleware", "middleware")
	defer span.End()

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

	_, err := repositories.GetUserSession(spanCtx, auth)
	if err != nil {
		log.Println("failed to get user session", err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	claims, err := jwt.ValidateToken(spanCtx, auth)
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

func MiddlewareRefreshToken(ctx *fiber.Ctx) error {

	span, spanCtx := apm.StartSpan(ctx.Context(), "MiddlewareRefreshToken", "middleware")
	defer span.End()

	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	var auth string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		auth = authHeader[7:]
	} else {
		auth = authHeader
	}

	// Validate refresh token exists in a database
	session, err := repositories.GetUserSessionByRefreshToken(spanCtx, auth)
	if err != nil {
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Invalid refresh token", nil)
	}

	// Check if the refresh token is expired in a database
	if time.Now().After(session.RefreshTokenExpired) {
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Refresh token expired", nil)
	}

	// Validate JWT structure (but allow expired tokens)
	claims, err := jwt.ValidateToken(spanCtx, auth)
	if err != nil {
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Invalid token format", nil)
	}

	ctx.Set("username", claims.Username)
	ctx.Set("full_name", claims.FullName)
	ctx.Set("refresh_token", auth)
	return ctx.Next()
}
