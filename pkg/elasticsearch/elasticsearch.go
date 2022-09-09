package elasticsearch

import (
	"emperror.dev/errors"
	v7 "github.com/olivere/elastic/v7"
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
		return nil, errors.WrapIf(err, "v7.NewClient")
	}

	return client, nil
}
