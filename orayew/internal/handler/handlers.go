package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Pong(ctx *fiber.Ctx) error {
	return successResponse(ctx, "pong")
}
