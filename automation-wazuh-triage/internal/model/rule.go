package model

import "automation-wazuh-triage/internal/entity"

// RuleResponse represents the simplified rule response for our API
type RuleResponse struct {
	Filename        string                  `json:"filename"`
	RelativeDirname string                  `json:"relative_dirname"`
	ID              int                     `json:"id"`
	Level           int                     `json:"level"`
	Status          string                  `json:"status"`
	Details         entity.WazuhRuleDetails `json:"details"`
	PciDss          []string                `json:"pci_dss"`
	Gpg13           []string                `json:"gpg13"`
	Gdpr            []string                `json:"gdpr"`
	Hipaa           []string                `json:"hipaa"`
	Nist80053       []string                `json:"nist_800_53"`
	Tsc             []string                `json:"tsc"`
	Mitre           []string                `json:"mitre"`
	Groups          []string                `json:"groups"`
	Description     string                  `json:"description"`
}

// ConvertWazuhRuleToResponse converts entity.WazuhRule to model.RuleResponse
func ConvertWazuhRuleToResponse(wazuhRule *entity.WazuhRule) *RuleResponse {
	return &RuleResponse{
		Filename:        wazuhRule.Filename,
		RelativeDirname: wazuhRule.RelativeDirname,
		ID:              wazuhRule.ID,
		Level:           wazuhRule.Level,
		Status:          wazuhRule.Status,
		Details:         wazuhRule.Details,
		PciDss:          wazuhRule.PciDss,
		Gpg13:           wazuhRule.Gpg13,
		Gdpr:            wazuhRule.Gdpr,
		Hipaa:           wazuhRule.Hipaa,
		Nist80053:       wazuhRule.Nist80053,
		Tsc:             wazuhRule.Tsc,
		Mitre:           wazuhRule.Mitre,
		Groups:          wazuhRule.Groups,
		Description:     wazuhRule.Description,
	}
}

// ConvertWazuhRulesToResponse converts slice of entity.WazuhRule to slice of model.RuleResponse
func ConvertWazuhRulesToResponse(wazuhRules []entity.WazuhRule) []RuleResponse {
	responses := make([]RuleResponse, len(wazuhRules))

	for i, wazuhRule := range wazuhRules {
		responses[i] = *ConvertWazuhRuleToResponse(&wazuhRule)
	}

	return responses
}

// ConvertWazuhAPIResponseToRules extracts rules from Wazuh API response and converts to model.RuleResponse
func ConvertWazuhAPIResponseToRules(apiResponse *entity.WazuhRulesAPIResponse) []RuleResponse {
	return ConvertWazuhRulesToResponse(apiResponse.Data.AffectedItems)
}
