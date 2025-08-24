package controllers

import (
	"go-chat-app/app/models"
	"go-chat-app/app/repositories"
	"go-chat-app/pkg/jwt"
	"go-chat-app/pkg/response"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx *fiber.Ctx) error {

	span, spanCtx := apm.StartSpan(ctx.Context(), "Register", "controller")
	defer span.End()

	user := new(models.User)

	if err := ctx.BodyParser(&user); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, "Invalid request format", err.Error())
	}

	if err := user.Validate(); err != nil {
		log.Printf("User validation failed: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
	}
	user.Password = string(hashPassword)

	err = repositories.CreateUser(spanCtx, user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Failed to create user", nil)
	}

	bodyResp := user.Username

	return response.SendSuccessResponse(ctx, bodyResp)
}

func LoginUser(ctx *fiber.Ctx) error {

	span, spanCtx := apm.StartSpan(ctx.Context(), "Login", "controller")
	defer span.End()

	now := time.Now()
	loginReq := new(models.LoginRequest)
	loginResp := new(models.LoginResponse)

	if err := ctx.BodyParser(&loginReq); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, "Invalid request format", err.Error())
	}

	if err := loginReq.Validate(); err != nil {
		log.Printf("User validation failed: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	user, err := repositories.GetUserByUsername(spanCtx, loginReq.Username)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "Validation failed", err.Error())
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		log.Printf("Failed to compare password: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Invalid credentials", err.Error())
	}

	token, err := jwt.GenerateToken(spanCtx, user.Username, user.FullName, `access`, now)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Internal server error", err.Error())
	}

	refreshToken, err := jwt.GenerateToken(spanCtx, user.Username, user.FullName, `refresh`, now)
	if err != nil {
		log.Printf("Failed to refresh token: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Internal server error", err.Error())
	}

	// Create a user session
	userSession := &models.UserSession{
		UserId:              user.Id,
		Token:               token,
		RefreshToken:        refreshToken,
		TokenExpired:        now.Add(jwt.MapTokenTypes[`access`]),
		RefreshTokenExpired: now.Add(jwt.MapTokenTypes[`refresh`]),
	}
	err = repositories.CreateUserSession(spanCtx, userSession)
	if err != nil {
		log.Printf("Failed to create user session: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Failed to create session", err.Error())
	}

	loginResp.Username = user.Username
	loginResp.FullName = user.FullName
	loginResp.Token = token
	loginResp.RefreshToken = refreshToken

	return response.SendSuccessResponse(ctx, loginResp)
}

func LogoutUser(ctx *fiber.Ctx) error {

	span, spanCtx := apm.StartSpan(ctx.Context(), "LogoutUser", "controller")
	defer span.End()

	authHeader := ctx.Get("Authorization")
	var token string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		token = authHeader
	}
	err := repositories.DeleteUserSession(spanCtx, token)
	if err != nil {
		log.Printf("Failed to delete user session: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Failed to delete session", err.Error())
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func RefreshToken(ctx *fiber.Ctx) error {

	span, spanCtx := apm.StartSpan(ctx.Context(), "RefreshToken", "controller")
	defer span.End()

	now := time.Now()
	username := ctx.Get("username")
	refreshToken := ctx.Get("refresh_token")
	fullName := ctx.Get("full_name")

	// Generate new tokens
	newAccessToken, err := jwt.GenerateToken(spanCtx, username, fullName, "access", now)
	if err != nil {
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Failed to generate access token", err.Error())
	}

	newRefreshToken, err := jwt.GenerateToken(spanCtx, username, fullName, "refresh", now)
	if err != nil {
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Failed to generate refresh token", err.Error())
	}

	// Update the session with new tokens and expiration times
	err = repositories.UpdateUserSessionTokens(spanCtx, newAccessToken, newRefreshToken,
		now.Add(jwt.MapTokenTypes["access"]), now.Add(jwt.MapTokenTypes["refresh"]), refreshToken)
	if err != nil {
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Failed to update session", err.Error())
	}

	return response.SendSuccessResponse(ctx, fiber.Map{
		"token":         newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
