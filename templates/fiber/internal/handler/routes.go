package handler

import "github.com/gofiber/fiber/v2"

func BindRoutes(app *fiber.App, h *Handler) {
	app.Get("/ping", h.Pong)
}
