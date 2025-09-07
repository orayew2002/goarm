package app

import (
	"github.com/gofiber/fiber/v2"

	"templates/internal/domain"
	"templates/internal/handler"
	"templates/internal/repo"
	"templates/internal/service"
)

func Run(appConfig domain.AppConfigs) error {
	repo := repo.NewRepo()
	service := service.NewService(repo)

	app := fiber.New()
	handler.BindRoutes(app, handler.NewHandler(service))

	return app.Listen(":" + appConfig.App.Port)
}
