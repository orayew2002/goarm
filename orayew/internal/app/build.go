package app

import (
	"github.com/gofiber/fiber/v2"

	"orayew/internal/domain"
	"orayew/internal/handler"
	"orayew/internal/repo"
	"orayew/internal/service"
	"orayew/pkg/pgxpool"
)

func Run(appConfig domain.AppConfigs) error {
	repo := repo.NewRepo(pgxpool.NewClient(appConfig.DB))
	service := service.NewService(repo)

	app := fiber.New()
	handler.BindRoutes(app, handler.NewHandler(service))

	return app.Listen(":" + appConfig.App.Port)
}
