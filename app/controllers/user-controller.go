package controllers

import (
	"go-chat-app/app/models"
	"go-chat-app/app/repositories"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx *fiber.Ctx) error {

	user := new(models.User)

	if err := ctx.BodyParser(&user); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
	}

	if err := user.Validate(); err != nil {
		log.Printf("User validation failed: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"error":   err.Error(),
		})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
	}
	user.Password = string(hashPassword)

	err = repositories.CreateUser(ctx.Context(), user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to create user",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"userId":  user.Id,
		"message": "User registered successfully",
	})
}
