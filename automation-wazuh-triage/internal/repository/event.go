package repository

import (
	"automation-wazuh-triage/internal/domain"
	"automation-wazuh-triage/pkg/logger"
	"context"

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
