package server

import (
	"context"
	"fmt"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/catalogs"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/web"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Log          logger.Logger
	Cfg          *config.Config
	EchoServer   customEcho.EchoHttpServer
	doneCh       chan struct{}
	GrpcServer   grpcServer.GrpcServer
	HealthServer *http.Server
}

func NewServer(log logger.Logger, cfg *config.Config) *Server {
	g := grpcServer.NewGrpcServer(cfg.GRPC, log)
	h := customEcho.NewEchoHttpServer(cfg.Http, log)
	return &Server{Log: log, Cfg: cfg, EchoServer: h, GrpcServer: g, HealthServer: NewHealthCheckServer(cfg)}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	catalogsConfigurator := catalogs.NewCatalogsServiceConfigurator(s.Log, s.Cfg, s.EchoServer, s.GrpcServer)
	err, catalogsCleanup := catalogsConfigurator.ConfigureCatalogsService(ctx)
	if err != nil {
		return err
	}
	defer catalogsCleanup()

	deliveryType := s.Cfg.DeliveryType

	s.RunMetrics(cancel)

	healthCleanup := s.RunHealthCheck(ctx)
	defer healthCleanup()

	switch deliveryType {
	case "http":
		go func() {
			if err := s.EchoServer.RunHttpServer(nil); err != nil {
				s.Log.Errorf("(s.RunHttpServer) err: {%v}", err)
				cancel()
			}
		}()
		s.Log.Infof("%s is listening on Http PORT: {%s}", web.GetMicroserviceName(s.Cfg), s.Cfg.Http.Port)

	case "grpc":
		go func() {
			if err := s.GrpcServer.RunGrpcServer(nil); err != nil {
				s.Log.Errorf("(s.RunGrpcServer) err: {%v}", err)
				cancel()
			}
		}()
		s.Log.Infof("%s is listening on Grpc PORT: {%s}", web.GetMicroserviceName(s.Cfg), s.Cfg.GRPC.Port)
	default:
		panic(fmt.Sprintf("server type %s is not supported", deliveryType))
	}

	<-ctx.Done()
	s.waitForShootDown(constants.WaitShotDownDuration)

	switch deliveryType {
	case "http":
		s.Log.Infof("%s is shutting down Http PORT: {%s}", web.GetMicroserviceName(s.Cfg), s.Cfg.Http.Port)
		if err := s.EchoServer.GracefulShutdown(ctx); err != nil {
			s.Log.Warnf("(Shutdown) err: {%v}", err)
		}
	case "grpc":
		s.Log.Infof("%s is shutting down Grpc PORT: {%s}", web.GetMicroserviceName(s.Cfg), s.Cfg.GRPC.Port)
		s.GrpcServer.GracefulShutdown()
	}

	<-s.doneCh
	s.Log.Infof("%s server exited properly", web.GetMicroserviceName(s.Cfg))

	return nil
}

func (s *Server) waitForShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		s.doneCh <- struct{}{}
	}()
}
