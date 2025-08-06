package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) Pong(ctx *gin.Context) {
	successResponse(ctx, "pong")
}
