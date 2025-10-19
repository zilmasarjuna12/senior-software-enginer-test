package repository

import (
	"automation-wazuh-triage/internal/domain"
	"automation-wazuh-triage/internal/entity"
	"automation-wazuh-triage/pkg/logger"
	"context"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/olivere/elastic/v7"
)

type wazuhEventRepository struct {
	openSearchClient *elastic.Client
}

func NewWazuhEventRepository(
	openSearchClient *elastic.Client,
) domain.WazuhEventRepository {
	return &wazuhEventRepository{
		openSearchClient: openSearchClient,
	}
}

func (r *wazuhEventRepository) FetchSecurityEvents(ctx context.Context, severity int, tags []string) ([]*elastic.SearchHit, error) {
	log := logger.WithRequestID(ctx)

	esQuery := elastic.NewBoolQuery().
		Must(
			elastic.NewRangeQuery("timestamp").
				Format("epoch_millis"),
		)

	searchResult, err := r.openSearchClient.Search().
		Index("wazuh-alerts-*").
		Sort("timestamp", false).
		Query(esQuery).
		Pretty(false).
		Do(context.Background())
	if err != nil {
		log.WithError(err).Error("[repository - event - GetSecurityEvents]: Failed to get security events`")
		return nil, err
	}

	return searchResult.Hits.Hits, nil
}

func (r *wazuhEventRepository) FetchSecurityEventByID(ctx context.Context, eventID string) (*entity.WazuhSecurityEvent, *elastic.SearchHit, error) {
	log := logger.WithRequestID(ctx)

	esQuery := elastic.NewBoolQuery().
		Filter(
			elastic.NewTermsQuery("id", eventID),
		)

	searchSource := elastic.NewSearchSource().
		Size(1).
		FetchSource(true).
		Query(esQuery)

	searchResult, err := r.openSearchClient.Search().
		Index("wazuh-alerts-*").
		SearchSource(searchSource).
		Do(context.Background())
	if err != nil {
		log.WithError(err).Error("[repository - event - FetchSecurityEventByID]: Failed to fetch security event by ID")
		return nil, nil, err
	}

	if len(searchResult.Hits.Hits) == 0 {
		return nil, nil, fmt.Errorf("event with ID %s not found", eventID)
	}

	var event entity.WazuhSecurityEvent
	if err := json.Unmarshal(searchResult.Hits.Hits[0].Source, &event); err != nil {
		log.WithError(err).Error("[repository - event - FetchSecurityEventByID]: Failed to unmarshal security event")
		return nil, nil, err
	}

	return &event, searchResult.Hits.Hits[0], nil
}
