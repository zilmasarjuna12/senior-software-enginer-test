package domain

import (
	"automation-wazuh-triage/internal/entity"
	"automation-wazuh-triage/internal/model"
	"context"

	"github.com/olivere/elastic/v7"
)

type WazuhEventRepository interface {
	FetchSecurityEvents(ctx context.Context, filter *model.FetchEventsRequest) (searchResults []*elastic.SearchHit, err error)
	FetchSecurityEventByID(ctx context.Context, eventID string) (event *entity.WazuhSecurityEvent, searchHit *elastic.SearchHit, err error)
}

type ClosedEventRepository interface {
	SaveClosedEvent(ctx context.Context, closedEvent *entity.ClosedEvent) error
	FetchClosedEvents(ctx context.Context) ([]*entity.ClosedEvent, error)
	FetchClosedEventByID(ctx context.Context, id string) (*entity.ClosedEvent, error)
	FetchClosedEventByEventID(ctx context.Context, eventID string) (*entity.ClosedEvent, error)
	UpdateClosedEventReason(ctx context.Context, id string, reason string) error
}

type EventUsecase interface {
	FetchEvents(ctx context.Context, filter *model.FetchEventsRequest) (searchResults []*elastic.SearchHit, err error)
	FetchEventsWithAutoClose(ctx context.Context, filter *model.FetchEventsRequest) (searchResults []*elastic.SearchHit, err error)
	AddEventToCloseEvent(ctx context.Context, eventID string, reason string) error
	FetchClosedEvents(ctx context.Context) ([]*entity.ClosedEvent, error)
	FetchClosedEventDetailsByID(ctx context.Context, id string) (*entity.ClosedEvent, *entity.WazuhRule, []entity.WazuhRule, error)
	UpdateClosedEventReason(ctx context.Context, id string, reason string) error
}
