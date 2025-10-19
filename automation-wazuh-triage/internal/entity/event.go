package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

type WazuhSecurityEventRule struct {
	Description string `json:"description"`
	Level       int    `json:"level"`
	ID          string `json:"id"`
}

type WazuhSecurityEvent struct {
	ID        json.RawMessage         `json:"id"`
	Timestamp string                  `json:"timestamp"`
	Rule      *WazuhSecurityEventRule `json:"rule"`
}

type ClosedEvent struct {
	ID       int       `json:"id" db:"id"`
	EventID  string    `json:"event_id" db:"event_id"`
	RuleID   string    `json:"rule_id" db:"rule_id"`
	RawEvent string    `json:"raw_event" db:"raw_event"`
	Reason   string    `json:"reason" db:"reason"`
	Status   string    `json:"status" db:"status"`
	CloseAt  time.Time `json:"close_at" db:"close_at"`
}

func (m *WazuhSecurityEvent) UnmarshalJSON(data []byte) error {
	type Alias WazuhSecurityEvent
	aux := &struct {
		ID interface{} `json:"id"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Convert ID to string
	switch v := aux.ID.(type) {
	case string:
		m.ID = json.RawMessage(v)
	case float64:
		m.ID = json.RawMessage(fmt.Sprintf("%.0f", v)) // convert number to string
	default:
		return fmt.Errorf("invalid ID type: %T", v)
	}

	return nil
}
