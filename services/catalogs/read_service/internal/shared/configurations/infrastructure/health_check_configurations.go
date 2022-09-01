package infrastructure

import (
	"context"
	"github.com/heptiolabs/healthcheck"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (ic *infrastructureConfigurator) configureHealthCheckEndpoints(ctx context.Context, mongoClient *mongo.Client) {

	health := healthcheck.NewHandler()

	health.AddReadinessCheck(constants.MongoDB, healthcheck.AsyncWithContext(ctx, func() error {
		if err := mongoClient.Ping(ctx, nil); err != nil {
			ic.log.Warnf("(MongoDB Readiness Check) err: {%v}", err)
			return err
		}
		return nil
	}, time.Duration(ic.cfg.Probes.CheckIntervalSeconds)*time.Second))

	//health.AddReadinessCheck(constants.ElasticSearch, healthcheck.AsyncWithContext(ctx, func() error {
	//	_, _, err := s.elasticClient.Ping(s.cfg.Elastic.URL).Do(ctx)
	//	if err != nil {
	//		s.log.Warnf("(ElasticSearch Readiness Check) err: {%v}", err)
	//		return http_errors.Wrap(err, "client.Ping")
	//	}
	//	return nil
	//}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))
	//
	//i.healthCheck.AddLivenessCheck(constants.ElasticSearch, healthcheck.AsyncWithContext(ctx, func() error {
	//	_, _, err := s.elasticClient.Ping(s.cfg.Elastic.URL).Do(ctx)
	//	if err != nil {
	//		s.log.Warnf("(ElasticSearch Liveness Check) err: {%v}", err)
	//		return http_errors.Wrap(err, "client.Ping")
	//	}
	//	return nil
	//}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))
}
