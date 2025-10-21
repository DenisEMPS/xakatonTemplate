package handler

import "github.com/gofiber/fiber/v2"

func (h *Handler) Template(c *fiber.Ctx) error {
	return c.SendString("hello")
}
