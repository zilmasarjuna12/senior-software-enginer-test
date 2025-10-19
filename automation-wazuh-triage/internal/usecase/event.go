package usecase

import (
	"automation-wazuh-triage/internal/domain"
	"automation-wazuh-triage/internal/entity"
	"automation-wazuh-triage/pkg/logger"
	"context"
	"encoding/json"
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

func (u *eventUsecase) FetchEvents(ctx context.Context, severity int, tags []string) (searchResults []*elastic.SearchHit, err error) {
	return u.wazuhEventRepo.FetchSecurityEvents(ctx, severity, tags)
}

func (u *eventUsecase) AddEventToCloseEvent(ctx context.Context, eventID string, reason string) error {
	log := logger.WithRequestID(ctx)

	securityEvent, resultElastic, err := u.wazuhEventRepo.FetchSecurityEventByID(ctx, eventID)
	if err != nil {
		log.WithError(err).Error("[usecase - event - CloseEvent]: Failed to fetch security event by ID")
		return err
	}

	// Convert elastic search result to JSON string
	resultElasticJSON, err := json.Marshal(resultElastic)
	if err != nil {
		log.WithError(err).Error("[usecase - event - CloseEvent]: Failed to marshal elastic search result to JSON")
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
