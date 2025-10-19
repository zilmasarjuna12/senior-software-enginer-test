package repository

import (
	"automation-wazuh-triage/internal/domain"
	"automation-wazuh-triage/internal/entity"
	"automation-wazuh-triage/pkg/logger"
	"automation-wazuh-triage/pkg/wazuh"
	"context"
	"encoding/json"
)

type ruleRepository struct {
}

func NewRuleRepository() domain.RuleRepository {
	return &ruleRepository{}
}

func (r *ruleRepository) GetDetailRules(ctx context.Context, ruleID string) (*entity.WazuhRule, error) {
	log := logger.WithRequestID(ctx)

	queryString := "rule_ids=" + ruleID

	client := wazuh.NewWazuh()

	responseBytes, err := client.GetRules(queryString)
	if err != nil {
		log.WithError(err).Error("[repository - rule - GetDetailRules]: Failed to get rule details")
		return nil, err
	}

	// Parse the Wazuh API response
	var apiResponse entity.WazuhRulesAPIResponse
	if err := json.Unmarshal(responseBytes, &apiResponse); err != nil {
		log.WithError(err).Error("[repository - rule - GetDetailRules]: Failed to unmarshal Wazuh API response")
		return nil, err
	}

	// Check if Wazuh API returned an error
	if apiResponse.Error != 0 {
		log.WithField("wazuh_error", apiResponse.Error).WithField("message", apiResponse.Message).Error("[repository - rule - GetDetailRules]: Wazuh API returned error")
		return nil, err
	}

	// Check if any rules were found
	if len(apiResponse.Data.AffectedItems) == 0 {
		log.WithField("rule_id", ruleID).Warn("[repository - rule - GetDetailRules]: No rule found with given ID")
		return nil, nil // Return nil to indicate not found
	}

	// Return the first (and should be only) rule
	rule := &apiResponse.Data.AffectedItems[0]
	log.WithField("rule_id", ruleID).Info("[repository - rule - GetDetailRules]: Successfully fetched rule details")

	return rule, nil
}

func (r *ruleRepository) GetListRulesByFiles(ctx context.Context, filename string) ([]entity.WazuhRule, error) {
	log := logger.WithRequestID(ctx)

	queryString := "filename=" + filename

	client := wazuh.NewWazuh()

	responseBytes, err := client.GetRules(queryString)
	if err != nil {
		log.WithError(err).Error("[repository - rule - GetListRulesByFiles]: Failed to get rules by file")
		return nil, err
	}

	// Parse the Wazuh API response
	var apiResponse entity.WazuhRulesAPIResponse
	if err := json.Unmarshal(responseBytes, &apiResponse); err != nil {
		log.WithError(err).Error("[repository - rule - GetListRulesByFiles]: Failed to unmarshal Wazuh API response")
		return nil, err
	}

	// Check if Wazuh API returned an error
	if apiResponse.Error != 0 {
		log.WithField("wazuh_error", apiResponse.Error).WithField("message", apiResponse.Message).Error("[repository - rule - GetListRulesByFiles]: Wazuh API returned error")
		return nil, err
	}

	log.WithField("filename", filename).WithField("rules_count", len(apiResponse.Data.AffectedItems)).Info("[repository - rule - GetListRulesByFiles]: Successfully fetched rules by file")

	return apiResponse.Data.AffectedItems, nil
}
