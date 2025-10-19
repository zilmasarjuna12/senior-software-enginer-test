package domain

import (
	"automation-wazuh-triage/internal/entity"
	"context"

	"github.com/olivere/elastic/v7"
)

type WazuhEventRepository interface {
	FetchSecurityEvents(ctx context.Context, severity int, tags []string) (searchResults []*elastic.SearchHit, err error)
	FetchSecurityEventByID(ctx context.Context, eventID string) (event *entity.WazuhSecurityEvent, searchHit *elastic.SearchHit, err error)
}

type ClosedEventRepository interface {
	SaveClosedEvent(ctx context.Context, closedEvent *entity.ClosedEvent) error
	FetchClosedEvents(ctx context.Context) ([]*entity.ClosedEvent, error)
	FetchClosedEventByID(ctx context.Context, id string) (*entity.ClosedEvent, error)
}

type EventUsecase interface {
	FetchEvents(ctx context.Context, severity int, tags []string) (searchResults []*elastic.SearchHit, err error)
	AddEventToCloseEvent(ctx context.Context, eventID string, reason string) error
	FetchClosedEvents(ctx context.Context) ([]*entity.ClosedEvent, error)
	FetchClosedEventDetailsByID(ctx context.Context, id string) (*entity.ClosedEvent, *entity.WazuhRule, []entity.WazuhRule, error)
}
