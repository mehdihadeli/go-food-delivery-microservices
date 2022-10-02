package server

import (
	"context"
	"fmt"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	webWoker "github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/catalogs"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/web/workers"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	log    logger.Logger
	cfg    *config.Config
	doneCh chan struct{}
}

func NewServer(log logger.Logger, cfg *config.Config) *Server {
	return &Server{log: log, cfg: cfg, doneCh: make(chan struct{})}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	grpcServer := grpcServer.NewGrpcServer(s.cfg.GRPC, s.log)
	echoServer := customEcho.NewEchoHttpServer(s.cfg.Http, s.log)

	ic := infrastructure.NewInfrastructureConfigurator(s.log, s.cfg)
	infrastructureConfigurations, err, infraCleanup := ic.ConfigInfrastructures(ctx)
	if err != nil {
		return err
	}
	defer infraCleanup()

	catalogsConfigurator := catalogs.NewCatalogsServiceConfigurator(infrastructureConfigurations, echoServer, grpcServer)
	err = catalogsConfigurator.ConfigureCatalogsService(ctx)
	if err != nil {
		return err
	}

	deliveryType := s.cfg.DeliveryType

	var serverError error

	switch deliveryType {
	case "http":
		go func() {
			if err := echoServer.RunHttpServer(ctx, nil); err != nil {
				s.log.Errorf("(s.RunHttpServer) err: {%v}", err)
				serverError = err
				cancel()
			}
		}()
		s.log.Infof("%s is listening on Http PORT: {%s}", s.cfg.GetMicroserviceNameUpper(), s.cfg.Http.Port)

	case "grpc":
		go func() {
			if err := grpcServer.RunGrpcServer(ctx, nil); err != nil {
				s.log.Errorf("(s.RunGrpcServer) err: {%v}", err)
				serverError = err
				cancel()
			}
		}()
		s.log.Infof("%s is listening on Grpc PORT: {%s}", s.cfg.GetMicroserviceNameUpper(), s.cfg.GRPC.Port)
	default:
		panic(fmt.Sprintf("server type %s is not supported", deliveryType))
	}

	backgroundWorkers := webWoker.NewWorkersRunner([]webWoker.Worker{
		workers.NewMetricsWorker(infrastructureConfigurations),
	})

	workersErr := backgroundWorkers.Start(ctx)
	go func() {
		for {
			select {
			case e := <-workersErr:
				serverError = e
				cancel()
				return
			}
		}
	}()

	// waiting for app get a canceled or completed signal
	<-ctx.Done()
	s.waitForShootDown(constants.WaitShotDownDuration)

	// waiting for shutdown time reached
	<-s.doneCh
	s.log.Infof("%s service exited properly", s.cfg.GetMicroserviceNameUpper())

	return serverError
}

func (s *Server) waitForShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		s.doneCh <- struct{}{}
	}()
}
