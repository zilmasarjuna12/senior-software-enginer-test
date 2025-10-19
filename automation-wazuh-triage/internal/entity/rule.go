package entity

// WazuhRuleDetails represents the details field in Wazuh rule response
type WazuhRuleDetails struct {
	NoAlert  string `json:"noalert,omitempty"`
	Category string `json:"category,omitempty"`
}

// WazuhRule represents a single rule from Wazuh API
type WazuhRule struct {
	Filename        string           `json:"filename"`
	RelativeDirname string           `json:"relative_dirname"`
	ID              int              `json:"id"`
	Level           int              `json:"level"`
	Status          string           `json:"status"`
	Details         WazuhRuleDetails `json:"details"`
	PciDss          []string         `json:"pci_dss"`
	Gpg13           []string         `json:"gpg13"`
	Gdpr            []string         `json:"gdpr"`
	Hipaa           []string         `json:"hipaa"`
	Nist80053       []string         `json:"nist_800_53"`
	Tsc             []string         `json:"tsc"`
	Mitre           []string         `json:"mitre"`
	Groups          []string         `json:"groups"`
	Description     string           `json:"description"`
}

// WazuhRulesAPIResponse represents the full Wazuh API response structure
type WazuhRulesAPIResponse struct {
	Data struct {
		AffectedItems      []WazuhRule `json:"affected_items"`
		TotalAffectedItems int         `json:"total_affected_items"`
		TotalFailedItems   int         `json:"total_failed_items"`
		FailedItems        []string    `json:"failed_items"`
	} `json:"data"`
	Message string `json:"message"`
	Error   int    `json:"error"`
}
