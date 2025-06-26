package handler

import "github.com/gin-gonic/gin"

func BindRoutes(app *gin.Engine, h *Handler) {
	app.GET("/ping", h.Pong)
}
