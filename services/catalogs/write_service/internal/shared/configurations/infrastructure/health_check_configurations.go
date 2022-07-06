package infrastructure

import (
	"context"
	"github.com/heptiolabs/healthcheck"
)

func (ic *infrastructureConfigurator) configureHealthCheckEndpoints(ctx context.Context) {

	health := healthcheck.NewHandler()
	ic.log.Info(health)

	//health.AddReadinessCheck(constants.ElasticSearch, healthcheck.AsyncWithContext(ctx, func() error {
	//	_, _, err := s.elasticClient.Ping(s.cfg.Elastic.URL).Do(ctx)
	//	if err != nil {
	//		s.log.Warnf("(ElasticSearch Readiness Check) err: {%v}", err)
	//		return errors.Wrap(err, "client.Ping")
	//	}
	//	return nil
	//}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))
	//
	//i.healthCheck.AddLivenessCheck(constants.ElasticSearch, healthcheck.AsyncWithContext(ctx, func() error {
	//	_, _, err := s.elasticClient.Ping(s.cfg.Elastic.URL).Do(ctx)
	//	if err != nil {
	//		s.log.Warnf("(ElasticSearch Liveness Check) err: {%v}", err)
	//		return errors.Wrap(err, "client.Ping")
	//	}
	//	return nil
	//}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))
}
