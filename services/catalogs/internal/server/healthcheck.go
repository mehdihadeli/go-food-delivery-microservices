package server

import (
	"context"
	"github.com/heptiolabs/healthcheck"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"net/http"
	"time"
)

func (s *Server) runHealthCheck(ctx context.Context) {
	health := healthcheck.NewHandler()

	mux := http.NewServeMux()
	s.hs = &http.Server{
		Handler:      mux,
		Addr:         s.cfg.Probes.Port,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}
	mux.HandleFunc(s.cfg.Probes.LivenessPath, health.LiveEndpoint)
	mux.HandleFunc(s.cfg.Probes.ReadinessPath, health.ReadyEndpoint)

	s.configureHealthCheckEndpoints(ctx, health)

	go func() {
		s.log.Infof("(%s) Kubernetes probes listening on port: {%s}", s.cfg.ServiceName, s.cfg.Probes.Port)
		if err := s.hs.ListenAndServe(); err != nil {
			s.log.Errorf("(ListenAndServe) err: {%v}", err)
		}
	}()
}

func (s *Server) configureHealthCheckEndpoints(ctx context.Context, health healthcheck.Handler) {

	health.AddReadinessCheck(constants.MongoDB, healthcheck.AsyncWithContext(ctx, func() error {
		if err := s.mongoClient.Ping(ctx, nil); err != nil {
			s.log.Warnf("(MongoDB Readiness Check) err: {%v}", err)
			return err
		}
		return nil
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	//health.AddReadinessCheck(constants.ElasticSearch, healthcheck.AsyncWithContext(ctx, func() error {
	//	_, _, err := s.elasticClient.Ping(s.cfg.Elastic.URL).Do(ctx)
	//	if err != nil {
	//		s.log.Warnf("(ElasticSearch Readiness Check) err: {%v}", err)
	//		return errors.Wrap(err, "client.Ping")
	//	}
	//	return nil
	//}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))
	//
	//health.AddLivenessCheck(constants.ElasticSearch, healthcheck.AsyncWithContext(ctx, func() error {
	//	_, _, err := s.elasticClient.Ping(s.cfg.Elastic.URL).Do(ctx)
	//	if err != nil {
	//		s.log.Warnf("(ElasticSearch Liveness Check) err: {%v}", err)
	//		return errors.Wrap(err, "client.Ping")
	//	}
	//	return nil
	//}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))
}

func (s *Server) shutDownHealthCheckServer(ctx context.Context) error {
	return s.hs.Shutdown(ctx)
}
