package controllers

import (
	"go-chat-app/app/models"
	"go-chat-app/app/repositories"
	"go-chat-app/pkg/response"
	"log"

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
		return response.SendFailureResponse(ctx, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	return response.SendSuccessResponse(ctx, fiber.Map{
		"userId":  user.Id,
		"message": "User registered successfully",
	})
}
