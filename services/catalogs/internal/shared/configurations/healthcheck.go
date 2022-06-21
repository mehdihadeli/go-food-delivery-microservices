package configurations

import (
	"context"
	"github.com/heptiolabs/healthcheck"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	web_constants "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/constants"
	"net/http"
	"time"
)

func (s *Server) runHealthCheck(ctx context.Context) {
	health := healthcheck.NewHandler()

	mux := http.NewServeMux()
	s.HttpServer = &http.Server{
		Handler:      mux,
		Addr:         s.Cfg.Probes.Port,
		WriteTimeout: web_constants.WriteTimeout,
		ReadTimeout:  web_constants.ReadTimeout,
	}
	mux.HandleFunc(s.Cfg.Probes.LivenessPath, health.LiveEndpoint)
	mux.HandleFunc(s.Cfg.Probes.ReadinessPath, health.ReadyEndpoint)

	s.configureHealthCheckEndpoints(ctx, health)

	go func() {
		s.Log.Infof("(%s) Kubernetes probes listening on port: {%s}", s.Cfg.ServiceName, s.Cfg.Probes.Port)
		if err := s.HttpServer.ListenAndServe(); err != nil {
			s.Log.Errorf("(ListenAndServe) err: {%v}", err)
		}
	}()
}

func (s *Server) configureHealthCheckEndpoints(ctx context.Context, health healthcheck.Handler) {

	health.AddReadinessCheck(constants.MongoDB, healthcheck.AsyncWithContext(ctx, func() error {
		if err := s.MongoClient.Ping(ctx, nil); err != nil {
			s.Log.Warnf("(MongoDB Readiness Check) err: {%v}", err)
			return err
		}
		return nil
	}, time.Duration(s.Cfg.Probes.CheckIntervalSeconds)*time.Second))

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
	return s.HttpServer.Shutdown(ctx)
}
