package server

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	webWoker "github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/catalogs"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/web/workers"
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

	ic := infrastructure.NewInfrastructureConfigurator(s.log, s.cfg)
	infrastructureConfigurations, infraCleanup, err := ic.ConfigInfrastructures(ctx)
	if err != nil {
		return err
	}
	defer infraCleanup()

	catalogsConfigurator := catalogs.NewCatalogsServiceConfigurator(infrastructureConfigurations)
	catalogConfigurations, err := catalogsConfigurator.ConfigureCatalogsService(ctx)
	if err != nil {
		return err
	}

	var serverError error

	go func() {
		if err := catalogConfigurations.CatalogsEchoServer.RunHttpServer(ctx, nil); err != nil {
			s.log.Errorf("(s.RunHttpServer) err: {%v}", err)
			serverError = err
			cancel()
		}
	}()
	s.log.Infof("%s is listening on Http PORT: {%s}", s.cfg.GetMicroserviceNameUpper(), s.cfg.Http.Port)

	go func() {
		if err := catalogConfigurations.CatalogsGrpcServer.RunGrpcServer(ctx, nil); err != nil {
			s.log.Errorf("(s.RunGrpcServer) err: {%v}", err)
			serverError = err
			cancel()
		}
	}()
	s.log.Infof("%s is listening on Grpc PORT: {%s}", s.cfg.GetMicroserviceNameUpper(), s.cfg.GRPC.Port)

	backgroundWorkers := webWoker.NewWorkersRunner([]webWoker.Worker{
		workers.NewRabbitMQWorker(s.log, catalogConfigurations.CatalogsBus),
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
	s.log.Infof("%s server exited properly", s.cfg.GetMicroserviceNameUpper())

	return serverError
}

func (s *Server) waitForShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		s.doneCh <- struct{}{}
	}()
}
