package middleware

import (
	"automation-wazuh-triage/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate a new request ID
		requestID := uuid.New().String()

		// Set request ID in context
		c.Locals("request_id", requestID)

		// Set request ID in response header
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}

// LoggingMiddleware logs incoming requests with request ID
func LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Create logger entry with request ID
		entry := logger.WithRequestID(c.Context())

		// Log request
		entry.WithFields(logrus.Fields{
			"method":     c.Method(),
			"path":       c.Path(),
			"ip":         c.IP(),
			"user_agent": c.Get("User-Agent"),
		}).Info("Request received")

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log response
		entry.WithFields(logrus.Fields{
			"method":   c.Method(),
			"path":     c.Path(),
			"status":   c.Response().StatusCode(),
			"duration": duration.String(),
			"size":     len(c.Response().Body()),
		}).Info("Request completed")

		return err
	}
}

// GetRequestID gets the request ID from fiber context
func GetRequestID(c *fiber.Ctx) string {
	if requestID, ok := c.Locals("request_id").(string); ok {
		return requestID
	}
	return ""
}
