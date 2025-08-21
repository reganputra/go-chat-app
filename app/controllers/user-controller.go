package controllers

import (
	"go-chat-app/app/models"
	"go-chat-app/app/repositories"
	"go-chat-app/pkg/jwt"
	"go-chat-app/pkg/response"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx *fiber.Ctx) error {

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

	err = repositories.CreateUser(ctx.Context(), user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Failed to create user", nil)
	}

	bodyResp := user.Username

	return response.SendSuccessResponse(ctx, bodyResp)
}

func LoginUser(ctx *fiber.Ctx) error {

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

	user, err := repositories.GetUserByUsername(ctx.Context(), loginReq.Username)
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

	token, err := jwt.GenerateToken(ctx.Context(), user.Username, user.FullName, `access`, now)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Internal server error", err.Error())
	}

	refreshToken, err := jwt.GenerateToken(ctx.Context(), user.Username, user.FullName, `refresh`, now)
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
	err = repositories.CreateUserSession(ctx.Context(), userSession)
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
	authHeader := ctx.Get("Authorization")
	var token string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		token = authHeader
	}
	err := repositories.DeleteUserSession(ctx.Context(), token)
	if err != nil {
		log.Printf("Failed to delete user session: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Failed to delete session", err.Error())
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func RefreshToken(ctx *fiber.Ctx) error {

	now := time.Now()
	username := ctx.Get("username")
	refreshToken := ctx.Get("Authorization")
	fullName := ctx.Get("full_name")

	token, err := jwt.GenerateToken(ctx.Context(), username, fullName, `access`, now)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Internal server error", err.Error())
	}

	err = repositories.UpdateUserSession(ctx.Context(), token, refreshToken)
	if err != nil {
		log.Printf("Failed to update user session: %v", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Failed to update session", err.Error())
	}

	return response.SendSuccessResponse(ctx, fiber.Map{
		"token": token,
	})
}
