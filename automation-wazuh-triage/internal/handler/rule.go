package handler

import (
	"automation-wazuh-triage/internal/domain"
	"automation-wazuh-triage/internal/model"
	"automation-wazuh-triage/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type RuleHandler struct {
	ruleUsecase domain.RuleUsecase
}

func NewRuleHandler(ruleUsecase domain.RuleUsecase) *RuleHandler {
	return &RuleHandler{
		ruleUsecase: ruleUsecase,
	}
}

func (h *RuleHandler) GetDetailRules(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	ruleID := c.Params("id")
	if ruleID == "" {
		log.Error("[handler]: Missing rule ID parameter")
		return c.Status(fiber.StatusBadRequest).JSON(model.NewResponseError("Missing rule ID parameter"))
	}

	// Get rule from usecase
	wazuhRule, err := h.ruleUsecase.GetDetailRules(c.Context(), ruleID)
	if err != nil {
		log.WithError(err).Error("[handler]: Failed to get rule details")
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewResponseError("Failed to get rule details"))
	}

	// Check if rule was found
	if wazuhRule == nil {
		log.WithField("rule_id", ruleID).Warn("[handler]: Rule not found")
		return c.Status(fiber.StatusNotFound).JSON(model.NewResponseError("Rule not found"))
	}

	// Convert to response format
	responseRule := model.ConvertWazuhRuleToResponse(wazuhRule)

	return c.Status(fiber.StatusOK).JSON(model.NewResponseSuccess(responseRule))
}

func (h *RuleHandler) GetListRulesByFiles(c *fiber.Ctx) error {
	log := logger.WithRequestID(c.Context())

	filename := c.Params("filename")
	if filename == "" {
		log.Error("[handler]: Missing filename parameter")
		return c.Status(fiber.StatusBadRequest).JSON(model.NewResponseError("Missing filename parameter"))
	}

	// Get rules from usecase
	wazuhRules, err := h.ruleUsecase.GetListRulesByFiles(c.Context(), filename)
	if err != nil {
		log.WithError(err).Error("[handler]: Failed to get rules by files")
		return c.Status(fiber.StatusInternalServerError).JSON(model.NewResponseError("Failed to get rules"))
	}

	// Convert to response format
	responseRules := model.ConvertWazuhRulesToResponse(wazuhRules)

	return c.Status(fiber.StatusOK).JSON(model.NewResponseSuccess(responseRules))
}
