package usecase

import (
	"automation-wazuh-triage/internal/domain"
	"context"

	"github.com/olivere/elastic/v7"
)

type eventUsecase struct {
	wazuhEventRepo domain.WazuhEventRepository
}

func NewEventUsecase(
	wazuhEventRepo domain.WazuhEventRepository,
) domain.EventUsecase {
	return &eventUsecase{
		wazuhEventRepo: wazuhEventRepo,
	}
}

func (u *eventUsecase) FetchEvents(ctx context.Context, severity int, tags []string) (searchResults []*elastic.SearchHit, err error) {
	return u.wazuhEventRepo.FetchSecurityEvents(ctx, severity, tags)
}
