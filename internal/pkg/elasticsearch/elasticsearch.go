package elasticsearch

import (
	"emperror.dev/errors"
	"github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	URL string `mapstructure:"url"`
}

func NewElasticClient(cfg Config) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.URL},
	})
	if err != nil {
		return nil, errors.WrapIf(err, "v8.elasticsearch")
	}

	return es, nil
}
