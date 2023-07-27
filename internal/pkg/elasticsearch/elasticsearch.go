package elasticsearch

import (
	"emperror.dev/errors"
	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticClient(cfg *ElasticOptions) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.URL},
	})
	if err != nil {
		return nil, errors.WrapIf(err, "v8.elasticsearch")
	}

	return es, nil
}
