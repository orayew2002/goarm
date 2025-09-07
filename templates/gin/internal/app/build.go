package app

import (
	"github.com/gin-gonic/gin"

	"templates/internal/domain"
	"templates/internal/handler"
	"templates/internal/repo"
	"templates/internal/service"
)

func Run(appConfig domain.AppConfigs) error {
	repo := repo.NewRepo()
	service := service.NewService(repo)

	app := gin.Default()
	handler.BindRoutes(app, handler.NewHandler(service))

	return app.Run(":" + appConfig.App.Port)
}
