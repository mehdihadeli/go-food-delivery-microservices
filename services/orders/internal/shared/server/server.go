package server

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	webWoker "github.com/mehdihadeli/store-golang-microservice-sample/pkg/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/orders"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web/workers"
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

	ordersConfigurator := orders.NewOrdersServiceConfigurator(infrastructureConfigurations)
	ordersConfigurations, err := ordersConfigurator.ConfigureOrdersService(ctx)
	if err != nil {
		return err
	}

	var serverError error

	go func() {
		if err := ordersConfigurations.OrdersEchoServer().RunHttpServer(ctx, nil); err != nil {
			s.log.Errorf("(s.RunHttpServer) err: {%v}", err)
			serverError = err
			cancel()
		}
	}()
	s.log.Infof("%s is listening on Http PORT: {%s}", web.GetMicroserviceName(s.cfg), s.cfg.Http.Port)

	go func() {
		if err := ordersConfigurations.OrdersGrpcServer().RunGrpcServer(ctx, nil); err != nil {
			s.log.Errorf("(s.RunGrpcServer) err: {%v}", err)
			serverError = err
			cancel()
		}
	}()
	s.log.Infof("%s is listening on Grpc PORT: {%s}", web.GetMicroserviceName(s.cfg), s.cfg.GRPC.Port)

	backgroundWorkers := webWoker.NewWorkersRunner([]webWoker.Worker{
		workers.NewRabbitMQWorker(s.log, ordersConfigurations.OrdersBus()), workers.NewEventStoreDBWorker(s.log, s.cfg, ordersConfigurations.OrdersSubscriptionAllWorker()),
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
	s.log.Infof("microservice %s exited successfully", web.GetMicroserviceName(s.cfg))

	return serverError
}

func (s *Server) waitForShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		s.doneCh <- struct{}{}
	}()
}
