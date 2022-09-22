package server

import (
	"context"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	rabbitmqBus "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/orders"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web"
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
	return &Server{log: log, cfg: cfg}
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

	ordersConfigurator := orders.NewOrdersServiceConfigurator(infrastructureConfigurations, echoServer, grpcServer)
	err = ordersConfigurator.ConfigureOrdersService(ctx)
	if err != nil {
		return err
	}

	deliveryType := s.cfg.DeliveryType

	s.RunMetrics(cancel)

	var serverError error

	esdbWorker := eventstroredb.NewEsdbSubscriptionAllWorker(
		s.log,
		infrastructureConfigurations.Esdb,
		s.cfg.EventStoreConfig,
		infrastructureConfigurations.EsdbSerializer,
		infrastructureConfigurations.CheckpointRepository,
		es.NewProjectionPublisher(infrastructureConfigurations.Projections))

	rabbitMQBus := rabbitmqBus.NewRabbitMQBus(s.log, infrastructureConfigurations.Consumers)
	defer rabbitMQBus.Stop(ctx)

	switch deliveryType {
	case "http":
		go func() {
			if err := echoServer.RunHttpServer(nil); err != nil {
				s.log.Errorf("(s.RunHttpServer) err: {%v}", err)
				serverError = err
				cancel()
			}
		}()
		s.log.Infof("%s is listening on Http PORT: {%s}", web.GetMicroserviceName(s.cfg), s.cfg.Http.Port)

	case "grpc":
		go func() {
			if err := grpcServer.RunGrpcServer(nil); err != nil {
				s.log.Errorf("(s.RunGrpcServer) err: {%v}", err)
				serverError = err
				cancel()
			}
		}()
		s.log.Infof("%s is listening on Grpc PORT: {%s}", web.GetMicroserviceName(s.cfg), s.cfg.GRPC.Port)
	default:
		panic(fmt.Sprintf("server type %s is not supported", deliveryType))
	}

	go func() {
		err := rabbitMQBus.Start(ctx)
		if err != nil {
			serverError = err
			cancel()
		}
	}()

	go func() {
		//https://developers.eventstore.com/clients/grpc/subscriptions.html#filtering-by-prefix-1
		option := &eventstroredb.EventStoreDBSubscriptionToAllOptions{
			FilterOptions: &esdb.SubscriptionFilter{
				Type:     esdb.StreamFilterType,
				Prefixes: s.cfg.Subscriptions.OrderSubscription.Prefix,
			},
			SubscriptionId: s.cfg.Subscriptions.OrderSubscription.SubscriptionId,
		}
		err := esdbWorker.SubscribeAll(ctx, option)
		if err != nil {
			s.log.Errorf("(esdbSubscribeAllWorker.SubscribeAll) err: {%v}", err)
			serverError = err
			cancel()
		}
	}()

	<-ctx.Done()
	s.waitForShootDown(constants.WaitShotDownDuration)

	switch deliveryType {
	case "http":
		s.log.Infof("%s is shutting down Http PORT: {%s}", web.GetMicroserviceName(s.cfg), s.cfg.Http.Port)
		if err := echoServer.GracefulShutdown(ctx); err != nil {
			s.log.Warnf("(Shutdown) err: {%v}", err)
		}
	case "grpc":
		s.log.Infof("%s is shutting down Grpc PORT: {%s}", web.GetMicroserviceName(s.cfg), s.cfg.GRPC.Port)
		grpcServer.GracefulShutdown()
	}

	<-s.doneCh
	s.log.Infof("%s server exited properly", web.GetMicroserviceName(s.cfg))

	return serverError
}

func (s *Server) waitForShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		s.doneCh <- struct{}{}
	}()
}
