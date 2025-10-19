package handler

import (
	"automation-wazuh-triage/internal/domain"
	"automation-wazuh-triage/internal/model"
	"automation-wazuh-triage/pkg/logger"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/olivere/elastic/v7"
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

	var req *model.FetchEventsRequest
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("[handler]: Failed to parse fetch events request")
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
	}

	var events []*elastic.SearchHit
	var err error

	// Check if auto-close feature is requested
	if req.AutoAddToClose {
		log.WithField("auto_add_to_close", true).Info("[handler]: Processing fetch events with auto-close enabled")
		events, err = h.eventUsecase.FetchEventsWithAutoClose(c.Context(), req)
	} else {
		events, err = h.eventUsecase.FetchEvents(c.Context(), req)
	}

	if err != nil {
		log.WithError(err).Error("[handler]: Failed to fetch events")
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch events")
	}

	// Prepare response with additional metadata if auto-close was used
	responseData := map[string]interface{}{
		"events": events,
	}

	if req.AutoAddToClose {
		responseData["auto_closed"] = true
		responseData["total_events_processed"] = len(events)

		// Add informative message
		message := "Events fetched and automatically added to closed events database"
		responseData["message"] = message
	}

	return c.Status(fiber.StatusOK).JSON(model.NewResponseSuccess(responseData))
}

func (h *EventHandler) AddToClose(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	// Get event_id from URL parameter
	eventID := c.Params("event_id")
	if eventID == "" {
		log.Error("[handler]: Missing event_id parameter")
		return c.Status(fiber.StatusBadRequest).JSON(model.NewResponseError("Missing event_id parameter"))
	}

	// Parse request body
	var req model.CloseEventRequest
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("[handler]: Failed to parse close event request")
		return c.Status(fiber.StatusBadRequest).JSON(model.NewResponseError("Invalid request payload"))
	}

	// Validate reason is provided
	if req.Reason == "" {
		log.Error("[handler]: Missing reason in request body")
		return c.Status(fiber.StatusBadRequest).JSON(model.NewResponseError("Reason is required"))
	}

	// Close the event
	err := h.eventUsecase.AddEventToCloseEvent(c.Context(), eventID, req.Reason)
	if err != nil {
		// Check if it's a duplicate event error
		if strings.Contains(err.Error(), "is already closed") {
			log.WithError(err).WithField("event_id", eventID).Warn("[handler]: Event already closed")
			return c.Status(fiber.StatusConflict).JSON(model.NewResponseError(err.Error()))
		}

		// Check if it's an event not found error
		if strings.Contains(err.Error(), "not found") {
			log.WithError(err).WithField("event_id", eventID).Warn("[handler]: Event not found")
			return c.Status(fiber.StatusNotFound).JSON(model.NewResponseError("Event not found"))
		}

		log.WithError(err).Error("[handler]: Failed to close event")
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewResponseError("Failed to close event"))
	}

	return c.Status(fiber.StatusOK).JSON(model.NewResponseSuccess(map[string]interface{}{
		"event_id": eventID,
		"status":   "closed",
		"message":  "Event successfully closed",
	}))
}

func (h *EventHandler) FetchClosedEvents(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	closedEvents, err := h.eventUsecase.FetchClosedEvents(c.Context())
	if err != nil {
		log.WithError(err).Error("[handler]: Failed to fetch closed events")
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewResponseError("Failed to fetch closed events"))
	}

	// Convert closed events to response format with parsed JSON
	responseEvents, err := model.ConvertClosedEventsToResponse(closedEvents)
	if err != nil {
		log.WithError(err).Error("[handler]: Failed to convert closed events to response format")
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewResponseError("Failed to process closed events"))
	}

	return c.Status(fiber.StatusOK).JSON(model.NewResponseSuccess(responseEvents))
}

func (h *EventHandler) FetchClosedEventByID(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	id := c.Params("id")
	if id == "" {
		log.Error("[handler]: Missing closed event ID parameter")
		return c.Status(fiber.StatusBadRequest).JSON(model.NewResponseError("Missing closed event ID parameter"))
	}

	// Get closed event with rule details
	closedEvent, ruleDetail, relatedRules, err := h.eventUsecase.FetchClosedEventDetailsByID(c.Context(), id)
	if err != nil {
		log.WithError(err).Error("[handler]: Failed to fetch closed event details by ID")
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewResponseError("Failed to fetch closed event details"))
	}

	// Check if closed event was found
	if closedEvent == nil {
		log.WithField("id", id).Warn("[handler]: Closed event not found")
		return c.Status(fiber.StatusNotFound).JSON(model.NewResponseError("Closed event not found"))
	}

	// Convert to detailed response format
	responseEvent, err := model.ConvertClosedEventToDetailResponse(closedEvent, ruleDetail, relatedRules)
	if err != nil {
		log.WithError(err).Error("[handler]: Failed to convert closed event to detail response format")
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewResponseError("Failed to process closed event details"))
	}

	return c.Status(fiber.StatusOK).JSON(model.NewResponseSuccess(responseEvent))
}

func (h *EventHandler) UpdateClosedEventReason(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	// Get closed event ID from URL parameter
	id := c.Params("id")
	if id == "" {
		log.Error("[handler]: Missing closed event ID parameter")
		return c.Status(fiber.StatusBadRequest).JSON(model.NewResponseError("Missing closed event ID parameter"))
	}

	// Parse request body
	var req model.UpdateClosedEventReasonRequest
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("[handler]: Failed to parse update reason request")
		return c.Status(fiber.StatusBadRequest).JSON(model.NewResponseError("Invalid request payload"))
	}

	// Validate reason is provided
	if req.Reason == "" {
		log.Error("[handler]: Missing reason in request body")
		return c.Status(fiber.StatusBadRequest).JSON(model.NewResponseError("Reason is required"))
	}

	// Update the closed event reason
	err := h.eventUsecase.UpdateClosedEventReason(c.Context(), id, req.Reason)
	if err != nil {
		// Check if it's a not found error
		if err.Error() == "closed event with ID "+id+" not found" {
			log.WithError(err).WithField("id", id).Warn("[handler]: Closed event not found")
			return c.Status(fiber.StatusNotFound).JSON(model.NewResponseError("Closed event not found"))
		}

		log.WithError(err).Error("[handler]: Failed to update closed event reason")
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewResponseError("Failed to update closed event reason"))
	}

	return c.Status(fiber.StatusOK).JSON(model.NewResponseSuccess(map[string]interface{}{
		"id":      id,
		"reason":  req.Reason,
		"message": "Closed event reason updated successfully",
	}))
}
