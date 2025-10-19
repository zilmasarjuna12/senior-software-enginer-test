package model

import (
	"automation-wazuh-triage/internal/entity"
	"encoding/json"
	"time"
)

type FetchEventsRequest struct {
	LevelRange     *RangeQuery `json:"level_range,omitempty"`
	Limit          int         `json:"limit,omitempty"`
	AutoAddToClose bool        `json:"auto_add_to_close,omitempty"`
	CloseReason    string      `json:"close_reason,omitempty"`
}

type RangeQuery struct {
	Gte interface{} `json:"gte,omitempty"` // Greater than or equal
	Gt  interface{} `json:"gt,omitempty"`  // Greater than
	Lte interface{} `json:"lte,omitempty"` // Less than or equal
	Lt  interface{} `json:"lt,omitempty"`  // Less than
}

type CloseEventRequest struct {
	Reason string `json:"reason"`
}

type UpdateClosedEventReasonRequest struct {
	Reason string `json:"reason"`
}

type ClosedEventResponse struct {
	ID       int         `json:"id"`
	EventID  string      `json:"event_id"`
	RuleID   string      `json:"rule_id"`
	RawEvent interface{} `json:"raw_event"` // This will hold the parsed JSON
	Reason   string      `json:"reason"`
	Status   string      `json:"status"`
	CloseAt  time.Time   `json:"close_at"`
}

type ClosedEventDetailResponse struct {
	ID           int            `json:"id"`
	EventID      string         `json:"event_id"`
	RuleID       string         `json:"rule_id"`
	RawEvent     interface{}    `json:"raw_event"` // This will hold the parsed JSON
	Reason       string         `json:"reason"`
	Status       string         `json:"status"`
	CloseAt      time.Time      `json:"close_at"`
	Rule         *RuleResponse  `json:"rule,omitempty"`          // Rule detail
	RuleAffected []RuleResponse `json:"rule_affected,omitempty"` // Related rules from same file
}

// ConvertClosedEventToResponse converts entity.ClosedEvent to model.ClosedEventResponse
// and parses the raw_event string into JSON object
func ConvertClosedEventToResponse(closedEvent *entity.ClosedEvent) (*ClosedEventResponse, error) {
	response := &ClosedEventResponse{
		ID:      closedEvent.ID,
		EventID: closedEvent.EventID,
		RuleID:  closedEvent.RuleID,
		Reason:  closedEvent.Reason,
		Status:  closedEvent.Status,
		CloseAt: closedEvent.CloseAt,
	}

	// Parse raw_event from JSON string to object
	if closedEvent.RawEvent != "" {
		var rawEventJSON interface{}
		if err := json.Unmarshal([]byte(closedEvent.RawEvent), &rawEventJSON); err != nil {
			// If parsing fails, return the raw string as fallback
			response.RawEvent = closedEvent.RawEvent
		} else {
			response.RawEvent = rawEventJSON
		}
	} else {
		response.RawEvent = nil
	}

	return response, nil
}

// ConvertClosedEventsToResponse converts slice of entity.ClosedEvent to slice of model.ClosedEventResponse
func ConvertClosedEventsToResponse(closedEvents []*entity.ClosedEvent) ([]*ClosedEventResponse, error) {
	responses := make([]*ClosedEventResponse, len(closedEvents))

	for i, closedEvent := range closedEvents {
		response, err := ConvertClosedEventToResponse(closedEvent)
		if err != nil {
			return nil, err
		}
		responses[i] = response
	}

	return responses, nil
}

// ConvertClosedEventToDetailResponse converts entity.ClosedEvent to model.ClosedEventDetailResponse
// with extended rule information
func ConvertClosedEventToDetailResponse(closedEvent *entity.ClosedEvent, rule *entity.WazuhRule, relatedRules []entity.WazuhRule) (*ClosedEventDetailResponse, error) {
	response := &ClosedEventDetailResponse{
		ID:      closedEvent.ID,
		EventID: closedEvent.EventID,
		RuleID:  closedEvent.RuleID,
		Reason:  closedEvent.Reason,
		Status:  closedEvent.Status,
		CloseAt: closedEvent.CloseAt,
	}

	// Parse raw_event from JSON string to object
	if closedEvent.RawEvent != "" {
		var rawEventJSON interface{}
		if err := json.Unmarshal([]byte(closedEvent.RawEvent), &rawEventJSON); err != nil {
			// If parsing fails, return the raw string as fallback
			response.RawEvent = closedEvent.RawEvent
		} else {
			response.RawEvent = rawEventJSON
		}
	} else {
		response.RawEvent = nil
	}

	// Add rule detail if available
	if rule != nil {
		response.Rule = ConvertWazuhRuleToResponse(rule)
	}

	// Add related rules if available
	if len(relatedRules) > 0 {
		response.RuleAffected = ConvertWazuhRulesToResponse(relatedRules)
	}

	return response, nil
}
