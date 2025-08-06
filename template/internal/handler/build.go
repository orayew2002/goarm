package handler

import (
	"template/pkg/http"

	"github.com/gin-gonic/gin"
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

func response(ctx *gin.Context, statusCode int, data any) {
	ctx.JSON(statusCode, jsonResponse{
		Success: true,
		Data:    data,
	})
}

func successResponse(ctx *gin.Context, data any) {
	response(ctx, http.StatusOK, data)
}

func errResponse(ctx *gin.Context, statusCode int, error error) {
	ctx.JSON(statusCode, jsonResponse{
		Success: true,
		Error:   error.Error(),
	})
}

func errBadResponse(ctx *gin.Context, err error) {
	errResponse(ctx, http.StatusBadRequest, err)
}

func errUnauthorizedResponse(ctx *gin.Context, err error) {
	errResponse(ctx, http.StatusUnauthorized, err)
}

func errForbiddenResponse(ctx *gin.Context, err error) {
	errResponse(ctx, http.StatusForbidden, err)
}

func errNotFoundResponse(ctx *gin.Context, err error) {
	errResponse(ctx, http.StatusNotFound, err)
}

func errConflictResponse(ctx *gin.Context, err error) {
	errResponse(ctx, http.StatusConflict, err)
}

func errTooManyRequestsResponse(ctx *gin.Context, err error) {
	errResponse(ctx, http.StatusTooManyRequests, err)
}

func errInternalServerErrorResponse(ctx *gin.Context, err error) {
	errResponse(ctx, http.StatusInternalServerError, err)
}

func errNotImplementedResponse(ctx *gin.Context, err error) {
	errResponse(ctx, http.StatusNotImplemented, err)
}

func errServiceUnavailableResponse(ctx *gin.Context, err error) {
	errResponse(ctx, http.StatusServiceUnavailable, err)
}
