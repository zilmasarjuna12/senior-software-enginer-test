package domain

import (
	"context"

	"github.com/olivere/elastic/v7"
)

type WazuhEventRepository interface {
	FetchSecurityEvents(ctx context.Context, severity int, tags []string) (searchResults []*elastic.SearchHit, err error)
}

type EventUsecase interface {
	FetchEvents(ctx context.Context, severity int, tags []string) (searchResults []*elastic.SearchHit, err error)
}
