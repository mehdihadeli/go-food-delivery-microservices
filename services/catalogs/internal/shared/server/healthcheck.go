package server

import (
	"context"
	"github.com/heptiolabs/healthcheck"
	web_constants "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/constants"
	"net/http"
)

func (s *Server) RunHealthCheck(ctx context.Context) func() {
	health := healthcheck.NewHandler()
	s.healthCheck = health

	mux := http.NewServeMux()
	httpServer := &http.Server{
		Handler:      mux,
		Addr:         s.Cfg.Probes.Port,
		WriteTimeout: web_constants.WriteTimeout,
		ReadTimeout:  web_constants.ReadTimeout,
	}
	mux.HandleFunc(s.Cfg.Probes.LivenessPath, health.LiveEndpoint)
	mux.HandleFunc(s.Cfg.Probes.ReadinessPath, health.ReadyEndpoint)

	go func() {
		s.Log.Infof("(%s) Kubernetes probes listening on port: {%s}", s.Cfg.ServiceName, s.Cfg.Probes.Port)
		if err := httpServer.ListenAndServe(); err != nil {
			s.Log.Errorf("(ListenAndServe) err: {%v}", err)
		}
	}()

	return func() {
		_ = shutDownHealthCheckServer(httpServer, ctx)
	}
}

func shutDownHealthCheckServer(server *http.Server, ctx context.Context) error {
	return server.Shutdown(ctx)
}
