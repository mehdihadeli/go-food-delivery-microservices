package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/catalogs"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/web"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	Log          logger.Logger
	Cfg          *config.Config
	Echo         *echo.Echo
	DoneCh       chan struct{}
	GrpcServer   *grpc.Server
	HealthServer *http.Server
}

func NewServer(log logger.Logger, cfg *config.Config) *Server {
	return &Server{Log: log, Cfg: cfg, Echo: echo.New(), GrpcServer: NewGrpcServer(), HealthServer: NewHealthCheckServer(cfg)}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	catalogsConfigurator := catalogs.NewCatalogsServiceConfigurator(s.Log, s.Cfg, s.Echo, s.GrpcServer)
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
			if err := s.RunHttpServer(nil); err != nil {
				s.Log.Errorf("(s.RunHttpServer) err: {%v}", err)
				cancel()
			}
		}()
		s.Log.Infof("%s is listening on Http PORT: {%s}", web.GetMicroserviceName(s.Cfg), s.Cfg.Http.Port)

	case "grpc":
		go func() {
			if err := s.RunGrpcServer(nil); err != nil {
				s.Log.Errorf("(s.RunGrpcServer) err: {%v}", err)
				cancel()
			}
		}()
		s.Log.Infof("%s is listening on Grpc PORT: {%s}", web.GetMicroserviceName(s.Cfg), s.Cfg.GRPC.Port)
	default:
		fmt.Sprintf("server type %s is not supported", deliveryType)
		//panic()
	}

	<-ctx.Done()
	s.WaitShootDown(constants.WaitShotDownDuration)

	switch deliveryType {
	case "http":
		s.Log.Infof("%s is shutting down Http PORT: {%s}", web.GetMicroserviceName(s.Cfg), s.Cfg.Http.Port)
		if err := s.Echo.Shutdown(ctx); err != nil {
			s.Log.Warnf("(Shutdown) err: {%v}", err)
		}
	case "grpc":
		s.Log.Infof("%s is shutting down Grpc PORT: {%s}", web.GetMicroserviceName(s.Cfg), s.Cfg.GRPC.Port)
		s.GrpcServer.Stop()
		s.GrpcServer.GracefulStop()
	}

	<-s.DoneCh
	s.Log.Infof("%s server exited properly", web.GetMicroserviceName(s.Cfg))

	return nil
}
