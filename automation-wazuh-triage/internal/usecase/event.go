package usecase

import (
	"automation-wazuh-triage/internal/domain"
	"automation-wazuh-triage/internal/entity"
	"automation-wazuh-triage/internal/model"
	"automation-wazuh-triage/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
)

type eventUsecase struct {
	wazuhEventRepo  domain.WazuhEventRepository
	closedEventRepo domain.ClosedEventRepository
	ruleRepo        domain.RuleRepository
}

func NewEventUsecase(
	wazuhEventRepo domain.WazuhEventRepository,
	closedEventRepo domain.ClosedEventRepository,
	ruleRepo domain.RuleRepository,
) domain.EventUsecase {
	return &eventUsecase{
		wazuhEventRepo:  wazuhEventRepo,
		closedEventRepo: closedEventRepo,
		ruleRepo:        ruleRepo,
	}
}

func (u *eventUsecase) FetchEvents(ctx context.Context, filter *model.FetchEventsRequest) (searchResults []*elastic.SearchHit, err error) {
	return u.wazuhEventRepo.FetchSecurityEvents(ctx, filter)
}

func (u *eventUsecase) FetchEventsWithAutoClose(ctx context.Context, filter *model.FetchEventsRequest) (searchResults []*elastic.SearchHit, err error) {
	log := logger.WithRequestID(ctx)

	// First, fetch the events
	searchResults, err = u.wazuhEventRepo.FetchSecurityEvents(ctx, filter)
	if err != nil {
		log.WithError(err).Error("[usecase - event - FetchEventsWithAutoClose]: Failed to fetch security events")
		return nil, err
	}

	// If autoAddToClose is enabled, process each event
	if filter.AutoAddToClose {
		closeReason := filter.CloseReason
		if closeReason == "" {
			closeReason = "Auto-closed by fetch request"
		}

		log.WithField("event_count", len(searchResults)).Info("[usecase - event - FetchEventsWithAutoClose]: Auto-closing events")

		successCount := 0
		skipCount := 0

		for _, hit := range searchResults {
			// Extract event ID from the hit
			eventRawID := hit.Id
			if eventRawID == "" {
				log.WithField("hit_id", hit.Id).Warn("[usecase - event - FetchEventsWithAutoClose]: Skipping event with missing ID")
				skipCount++
				continue
			}

			// Parse the event to get rule ID
			var securityEvent entity.WazuhSecurityEvent
			if err := json.Unmarshal(hit.Source, &securityEvent); err != nil {
				log.WithError(err).WithField("event_id", eventRawID).Warn("[usecase - event - FetchEventsWithAutoClose]: Failed to parse event, skipping auto-close")
				skipCount++
				continue
			}

			eventID := string(securityEvent.ID)

			// Check if the event is already closed
			existingClosedEvent, err := u.closedEventRepo.FetchClosedEventByEventID(ctx, eventID)
			if err != nil {
				log.WithError(err).WithField("event_id", eventID).Warn("[usecase - event - FetchEventsWithAutoClose]: Failed to check existing closed event, skipping auto-close")
				skipCount++
				continue
			}

			if existingClosedEvent != nil {
				log.WithField("event_id", eventID).WithField("existing_closed_id", existingClosedEvent.ID).Debug("[usecase - event - FetchEventsWithAutoClose]: Event already closed, skipping")
				skipCount++
				continue
			}

			// Convert search hit to JSON string for storage
			hitJSON, err := json.Marshal(hit)
			if err != nil {
				log.WithError(err).WithField("event_id", eventID).Warn("[usecase - event - FetchEventsWithAutoClose]: Failed to marshal hit to JSON, skipping auto-close")
				skipCount++
				continue
			}

			// Create closed event record
			closedEvent := &entity.ClosedEvent{
				EventID:  eventID,
				RuleID:   securityEvent.Rule.ID,
				RawEvent: string(hitJSON),
				Reason:   closeReason,
				Status:   "closed",
				CloseAt:  time.Now(),
			}

			// Save to closed events database
			if err := u.closedEventRepo.SaveClosedEvent(ctx, closedEvent); err != nil {
				log.WithError(err).WithField("event_id", eventID).Error("[usecase - event - FetchEventsWithAutoClose]: Failed to save closed event, continuing with other events")
				skipCount++
				continue
			}

			log.WithField("event_id", eventID).Debug("[usecase - event - FetchEventsWithAutoClose]: Successfully auto-closed event")
			successCount++
		}

		log.WithField("processed_events", len(searchResults)).WithField("success_count", successCount).WithField("skip_count", skipCount).Info("[usecase - event - FetchEventsWithAutoClose]: Completed auto-closing process")
	}

	return searchResults, nil
}

func (u *eventUsecase) AddEventToCloseEvent(ctx context.Context, eventID string, reason string) error {
	log := logger.WithRequestID(ctx)

	// Check if the event is already closed
	existingClosedEvent, err := u.closedEventRepo.FetchClosedEventByEventID(ctx, eventID)
	if err != nil {
		log.WithError(err).Error("[usecase - event - AddEventToCloseEvent]: Failed to check existing closed event")
		return err
	}

	if existingClosedEvent != nil {
		log.WithField("event_id", eventID).WithField("existing_closed_id", existingClosedEvent.ID).Warn("[usecase - event - AddEventToCloseEvent]: Event already closed")
		return fmt.Errorf("event with ID %s is already closed (closed event ID: %d)", eventID, existingClosedEvent.ID)
	}

	securityEvent, resultElastic, err := u.wazuhEventRepo.FetchSecurityEventByID(ctx, eventID)
	if err != nil {
		log.WithError(err).Error("[usecase - event - AddEventToCloseEvent]: Failed to fetch security event by ID")
		return err
	}

	// Convert elastic search result to JSON string
	resultElasticJSON, err := json.Marshal(resultElastic)
	if err != nil {
		log.WithError(err).Error("[usecase - event - AddEventToCloseEvent]: Failed to marshal elastic search result to JSON")
		return err
	}

	closedEvent := &entity.ClosedEvent{
		EventID:  eventID,
		RuleID:   securityEvent.Rule.ID,     // This would be fetched from the actual event
		RawEvent: string(resultElasticJSON), // This would be the full event JSON
		Reason:   reason,
		Status:   "closed",
		CloseAt:  time.Now(),
	}

	return u.closedEventRepo.SaveClosedEvent(ctx, closedEvent)
}

func (u *eventUsecase) FetchClosedEvents(ctx context.Context) ([]*entity.ClosedEvent, error) {
	return u.closedEventRepo.FetchClosedEvents(ctx)
}

func (u *eventUsecase) FetchClosedEventDetailsByID(ctx context.Context, id string) (*entity.ClosedEvent, *entity.WazuhRule, []entity.WazuhRule, error) {
	log := logger.WithRequestID(ctx)

	// Get the closed event first
	closedEvent, err := u.closedEventRepo.FetchClosedEventByID(ctx, id)
	if err != nil {
		log.WithError(err).Error("[usecase - event - FetchClosedEventDetailsByID]: Failed to fetch closed event by ID")
		return nil, nil, nil, err
	}

	if closedEvent == nil {
		log.WithField("id", id).Warn("[usecase - event - FetchClosedEventDetailsByID]: Closed event not found")
		return nil, nil, nil, nil
	}

	// Get rule details if rule_id is available
	var ruleDetail *entity.WazuhRule
	var relatedRules []entity.WazuhRule

	if closedEvent.RuleID != "" {
		// Get the specific rule detail
		ruleDetail, err = u.ruleRepo.GetDetailRules(ctx, closedEvent.RuleID)
		if err != nil {
			log.WithError(err).WithField("rule_id", closedEvent.RuleID).Warn("[usecase - event - FetchClosedEventDetailsByID]: Failed to fetch rule details, continuing without rule info")
		}

		// If we got the rule detail and it has a filename, get related rules from the same file
		if ruleDetail != nil && ruleDetail.Filename != "" {
			relatedRules, err = u.ruleRepo.GetListRulesByFiles(ctx, ruleDetail.Filename)
			if err != nil {
				log.WithError(err).WithField("filename", ruleDetail.Filename).Warn("[usecase - event - FetchClosedEventDetailsByID]: Failed to fetch related rules, continuing without related rules")
				relatedRules = []entity.WazuhRule{} // Set empty slice instead of nil
			}
		}
	}

	return closedEvent, ruleDetail, relatedRules, nil
}

func (u *eventUsecase) UpdateClosedEventReason(ctx context.Context, id string, reason string) error {
	log := logger.WithRequestID(ctx)

	// Validate that the reason is not empty
	if reason == "" {
		log.Error("[usecase - event - UpdateClosedEventReason]: Reason cannot be empty")
		return fmt.Errorf("reason cannot be empty")
	}

	// Check if the closed event exists first
	closedEvent, err := u.closedEventRepo.FetchClosedEventByID(ctx, id)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("[usecase - event - UpdateClosedEventReason]: Failed to fetch closed event by ID")
		return err
	}

	if closedEvent == nil {
		log.WithField("id", id).Warn("[usecase - event - UpdateClosedEventReason]: Closed event not found")
		return fmt.Errorf("closed event with ID %s not found", id)
	}

	// Update the reason
	err = u.closedEventRepo.UpdateClosedEventReason(ctx, id, reason)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("[usecase - event - UpdateClosedEventReason]: Failed to update closed event reason")
		return err
	}

	log.WithField("id", id).WithField("reason", reason).Info("[usecase - event - UpdateClosedEventReason]: Successfully updated closed event reason")
	return nil
}
