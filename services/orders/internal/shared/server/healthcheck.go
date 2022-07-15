package server

import (
	"context"
	"github.com/heptiolabs/healthcheck"
	"net/http"
	"thub.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"thub.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/constants"
)

func NewHealthCheckServer(config *config.Config) *http.Server {
	mux := http.NewServeMux()
	health := healthcheck.NewHandler()
	mux.HandleFunc(config.Probes.LivenessPath, health.LiveEndpoint)
	mux.HandleFunc(config.Probes.ReadinessPath, health.ReadyEndpoint)

	httpServer := &http.Server{
		Handler:      mux,
		Addr:         config.Probes.Port,
		WriteTimeout: constants.WriteTimeout,
		ReadTimeout:  constants.ReadTimeout,
	}

	return httpServer
}

func (s *Server) RunHealthCheck(ctx context.Context) func() {

	go func() {
		s.Log.Infof("(%s) Kubernetes probes listening on port: {%s}", s.Cfg.ServiceName, s.Cfg.Probes.Port)
		if err := s.HealthServer.ListenAndServe(); err != nil {
			s.Log.Errorf("(ListenAndServe) err: {%v}", err)
		}
	}()

	return func() {
		_ = shutDownHealthCheckServer(s.HealthServer, ctx)
	}
}

func shutDownHealthCheckServer(server *http.Server, ctx context.Context) error {
	return server.Shutdown(ctx)
}
