package handler

import (
	"github.com/gofiber/fiber/v2"

	"orayew/pkg/http"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

type jsonResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Data    any    `json:"data"`
}

func response(ctx *fiber.Ctx, statusCode int, data any) error {
	return ctx.Status(statusCode).JSON(jsonResponse{
		Success: true,
		Data:    data,
	})
}

func successResponse(ctx *fiber.Ctx, data any) error {
	return response(ctx, http.StatusOK, data)
}

func errResponse(ctx *fiber.Ctx, statusCode int, error error) error {
	return ctx.Status(statusCode).JSON(jsonResponse{
		Success: false,
		Error:   error.Error(),
	})
}

func errBadResponse(ctx *fiber.Ctx, err error) error {
	return errResponse(ctx, http.StatusBadRequest, err)
}

func errUnauthorizedResponse(ctx *fiber.Ctx, err error) error {
	return errResponse(ctx, http.StatusUnauthorized, err)
}

func errForbiddenResponse(ctx *fiber.Ctx, err error) error {
	return errResponse(ctx, http.StatusForbidden, err)
}

func errNotFoundResponse(ctx *fiber.Ctx, err error) error {
	return errResponse(ctx, http.StatusNotFound, err)
}

func errConflictResponse(ctx *fiber.Ctx, err error) error {
	return errResponse(ctx, http.StatusConflict, err)
}

func errTooManyRequestsResponse(ctx *fiber.Ctx, err error) error {
	return errResponse(ctx, http.StatusTooManyRequests, err)
}

func errInternalServerErrorResponse(ctx *fiber.Ctx, err error) error {
	return errResponse(ctx, http.StatusInternalServerError, err)
}

func errNotImplementedResponse(ctx *fiber.Ctx, err error) error {
	return errResponse(ctx, http.StatusNotImplemented, err)
}

func errServiceUnavailableResponse(ctx *fiber.Ctx, err error) error {
	return errResponse(ctx, http.StatusServiceUnavailable, err)
}
