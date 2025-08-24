package controllers

import (
	"go-chat-app/app/repositories"
	"go-chat-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm"
)

func GetMessagesHistory(ctx *fiber.Ctx) error {

	span, spanCtx := apm.StartSpan(ctx.Context(), "GetMessagesHistory", "controller")
	defer span.End()

	resp, err := repositories.GetAllMessage(spanCtx)
	if err != nil {
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Internal Server Error", nil)
	}
	return response.SendSuccessResponse(ctx, resp)
}
