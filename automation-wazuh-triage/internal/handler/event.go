package handler

import (
	"automation-wazuh-triage/internal/domain"
	"automation-wazuh-triage/internal/model"
	"automation-wazuh-triage/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type EventHandler struct {
	eventUsecase domain.EventUsecase
}

func NewEventHandler(eventUsecase domain.EventUsecase) *EventHandler {
	return &EventHandler{
		eventUsecase: eventUsecase,
	}
}
func (h *EventHandler) FetchEvents(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	var req model.FetchEventsRequest
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("[handler]: Failed to parse fetch events request")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
	}

	events, err := h.eventUsecase.FetchEvents(c.Context(), req.Severity, req.Tags)

	log.Info("Fetching events")

	if err != nil {
		log.WithError(err).Error("[handler]: Failed to fetch events")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch events")
	}

	return c.Status(fiber.StatusOK).JSON(model.NewResponseSuccess(events))
}
