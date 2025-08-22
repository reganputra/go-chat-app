package controllers

import (
	"go-chat-app/app/repositories"
	"go-chat-app/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func GetMessagesHistory(ctx *fiber.Ctx) error {
	resp, err := repositories.GetAllMessage(ctx.Context())
	if err != nil {
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Internal Server Error", nil)
	}
	return response.SendSuccessResponse(ctx, resp)
}
