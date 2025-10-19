package main

import (
	"automation-wazuh-triage/internal/route"
	"automation-wazuh-triage/pkg/logger"
	"automation-wazuh-triage/pkg/opensearch"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger()
	log := logger.GetLogger()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found")
	}

	opensearch := opensearch.NewOpenSearch(log)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log := logger.WithRequestID(c.Context())

			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			log.WithError(err).WithField("status_code", code).Error("Request error")

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	route.SetupRoutes(app, opensearch)
	log.Fatal(app.Listen(":" + port))
}
