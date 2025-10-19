package route

import (
	"automation-wazuh-triage/internal/handler"
	"automation-wazuh-triage/internal/repository"
	"automation-wazuh-triage/internal/usecase"
	"automation-wazuh-triage/pkg/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/olivere/elastic/v7"
)

func SetupRoutes(app *fiber.App, openSearchClient *elastic.Client) {
	eventRepository := repository.NewWazuhEventRepository(openSearchClient)
	eventUsecase := usecase.NewEventUsecase(eventRepository)
	eventHandler := handler.NewEventHandler(eventUsecase)

	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.LoggingMiddleware())
	app.Use(recover.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success":   true,
			"message":   "success",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	v1 := app.Group("/v1")

	v1.Post("/events", eventHandler.FetchEvents)
}
