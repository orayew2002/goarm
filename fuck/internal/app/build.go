package app

import (
	"fuck/internal/domain"
	"fuck/internal/handler"
	"fuck/internal/repo"
	"fuck/internal/service"

	"github.com/gin-gonic/gin"
)

func Run(appConfig domain.AppConfigs) error {
	repo := repo.NewRepo()
	service := service.NewService(repo)

	app := gin.Default()
	handler.BindRoutes(app, handler.NewHandler(service))

	return app.Run(":" + appConfig.App.Port)
}
