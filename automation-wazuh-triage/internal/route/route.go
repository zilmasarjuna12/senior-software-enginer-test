package route

import (
	"automation-wazuh-triage/internal/handler"
	"automation-wazuh-triage/internal/repository"
	"automation-wazuh-triage/internal/usecase"
	"automation-wazuh-triage/pkg/database"
	"automation-wazuh-triage/pkg/middleware"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/olivere/elastic/v7"
)

func SetupRoutes(app *fiber.App, openSearchClient *elastic.Client) {
	// Initialize SQLite database
	db, err := database.InitSQLite("./data/events.db")
	if err != nil {
		log.Fatalf("Failed to initialize SQLite database: %v", err)
	}

	// Initialize repositories
	eventRepository := repository.NewWazuhEventRepository(openSearchClient)
	closedEventRepository := repository.NewClosedEventRepository(db)
	ruleRepository := repository.NewRuleRepository()

	// Initialize usecase
	eventUsecase := usecase.NewEventUsecase(eventRepository, closedEventRepository, ruleRepository)
	ruleUsecase := usecase.NewRuleUsecase(ruleRepository)

	// Initialize handler
	eventHandler := handler.NewEventHandler(eventUsecase)
	ruleHandler := handler.NewRuleHandler(ruleUsecase)

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

	// Swagger documentation
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/docs/openapi.yaml",
		DeepLinking: false,
	}))

	// Serve the OpenAPI specification file
	app.Static("/docs", "./docs")

	v1 := app.Group("/v1")

	v1.Post("/events", eventHandler.FetchEvents)
	v1.Post("/events/:event_id/close", eventHandler.AddToClose)
	v1.Get("/events/close", eventHandler.FetchClosedEvents)
	v1.Get("/events/close/:id", eventHandler.FetchClosedEventByID)
	v1.Patch("/events/close/:id/reason", eventHandler.UpdateClosedEventReason)

	v1.Get("/rules/:id", ruleHandler.GetDetailRules)
	v1.Get("/rules/file/:filename", ruleHandler.GetListRulesByFiles)
}
