package infrastructure

import (
	"context"
	"emperror.dev/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/elasticsearch"
	v7 "github.com/olivere/elastic/v7"
)

func (ic *infrastructureConfigurator) configElasticSearch(ctx context.Context) (*v7.Client, error, func()) {

	elasticClient, err := elasticsearch.NewElasticClient(ic.cfg.Elastic)
	if err != nil {
		return nil, err, nil
	}

	info, code, err := elasticClient.Ping(ic.cfg.Elastic.URL).Do(ctx)
	if err != nil {
		return nil, errors.WrapIf(err, "client.Ping"), nil
	}
	ic.log.Infof("Elasticsearch returned with code {%d} and version {%s}", code, info.Version.Number)

	esVersion, err := elasticClient.ElasticsearchVersion(ic.cfg.Elastic.URL)
	if err != nil {
		return nil, errors.WrapIf(err, "client.ElasticsearchVersion"), nil
	}
	ic.log.Infof("Elasticsearch version {%s}", esVersion)

	return elasticClient, nil, func() {}
}
