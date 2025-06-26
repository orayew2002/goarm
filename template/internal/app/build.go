package app

import (
	"template/internal/domain"
	"template/internal/handler"
	"template/internal/repo"
	"template/internal/service"

	"github.com/gin-gonic/gin"
)

func Run(appConfig domain.AppConfigs) error {
	repo := repo.NewRepo()
	service := service.NewService(repo)

	app := gin.Default()
	handler.BindRoutes(app, handler.NewHandler(service))

	return app.Run(":" + appConfig.App.Port)
}
