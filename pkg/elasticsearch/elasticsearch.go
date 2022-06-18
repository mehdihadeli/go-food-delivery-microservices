package elasticsearch

import (
	v7 "github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
)

type Config struct {
	URL         string `mapstructure:"url"`
	Sniff       bool   `mapstructure:"sniff"`
	Gzip        bool   `mapstructure:"gzip"`
	Explain     bool   `mapstructure:"explain"`
	FetchSource bool   `mapstructure:"fetchSource"`
	Version     bool   `mapstructure:"version"`
	Pretty      bool   `mapstructure:"pretty"`
}

func NewElasticClient(cfg Config) (*v7.Client, error) {
	client, err := v7.NewClient(
		v7.SetURL(cfg.URL),
		v7.SetSniff(cfg.Sniff),
		v7.SetGzip(cfg.Gzip),
	)
	if err != nil {
		return nil, errors.Wrap(err, "v7.NewClient")
	}

	return client, nil
}
