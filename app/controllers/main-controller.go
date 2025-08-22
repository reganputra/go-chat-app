package controllers

import "github.com/gofiber/fiber/v2"

func RenderUI(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func RenderAuth(c *fiber.Ctx) error {
	return c.Render("auth", fiber.Map{})
}

func RenderChat(c *fiber.Ctx) error {
	return c.Render("chat", fiber.Map{})
}
