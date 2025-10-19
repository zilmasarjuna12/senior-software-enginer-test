package opensearch

import (
	"os"

	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

func NewOpenSearch(
	log *logrus.Logger,
) *elastic.Client {
	host := os.Getenv("INDEXER_HOST")
	username := os.Getenv("INDEXER_USERNAME")
	password := os.Getenv("INDEXER_PASSWORD")

	client, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetBasicAuth(username, password),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}

	return client
}
