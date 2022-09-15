package e2e

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/projections"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
	"math"
	"net/http/httptest"
	"time"
)

type E2ETestFixture struct {
	Echo *echo.Echo
	*infrastructure.InfrastructureConfiguration
	V1         *V1Groups
	GrpcServer grpcServer.GrpcServer
	HttpServer *httptest.Server
	EsdbWorker eventstroredb.EsdbSubscriptionAllWorker
	Cleanup    func()
	ctx        context.Context
	cancel     context.CancelFunc
}

type V1Groups struct {
	OrdersGroup *echo.Group
}

func NewE2ETestFixture() *E2ETestFixture {
	deadline := time.Now().Add(time.Duration(math.MaxInt64))
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	cfg, _ := config.InitConfig("test")
	c := infrastructure.NewInfrastructureConfigurator(defaultLogger.Logger, cfg)
	infrastructures, _, cleanup := c.ConfigInfrastructures(ctx)
	echo := echo.New()

	v1Group := echo.Group("/api/v1")
	ordersV1 := v1Group.Group("/orders")

	v1Groups := &V1Groups{OrdersGroup: ordersV1}

	err := mediatr.ConfigOrdersMediator(infrastructures)
	if err != nil {
		cancel()
		return nil
	}

	projections.ConfigOrderProjections(infrastructures)

	err = mappings.ConfigureMappings()
	if err != nil {
		cancel()
		return nil
	}

	httpServer := httptest.NewServer(echo)
	grpcServer := grpcServer.NewGrpcServer(cfg.GRPC, defaultLogger.Logger)

	esdbSubscribeAllWorker := eventstroredb.NewEsdbSubscriptionAllWorker(
		infrastructures.Log,
		infrastructures.Esdb,
		infrastructures.Cfg.EventStoreConfig,
		infrastructures.EsdbSerializer,
		infrastructures.CheckpointRepository,
		es.NewProjectionPublisher(infrastructures.Projections))

	return &E2ETestFixture{
		Cleanup: func() {
			cleanup()
			grpcServer.GracefulShutdown()
			echo.Shutdown(ctx)
			httpServer.Close()
		},
		InfrastructureConfiguration: infrastructures,
		Echo:                        echo,
		V1:                          v1Groups,
		GrpcServer:                  grpcServer,
		HttpServer:                  httpServer,
		EsdbWorker:                  esdbSubscribeAllWorker,
		ctx:                         ctx,
		cancel:                      cancel,
	}
}

func (e *E2ETestFixture) Run() {
	go func() {
		if err := e.GrpcServer.RunGrpcServer(nil); err != nil {
			e.cancel()
			e.Log.Errorf("(s.RunGrpcServer) err: %v", err)
		}
	}()

	go func() {
		//https://developers.eventstore.com/clients/grpc/subscriptions.html#filtering-by-prefix-1
		option := &eventstroredb.EventStoreDBSubscriptionToAllOptions{
			FilterOptions: &esdb.SubscriptionFilter{
				Type:     esdb.StreamFilterType,
				Prefixes: e.Cfg.Subscriptions.OrderSubscription.Prefix,
			},
			SubscriptionId: e.Cfg.Subscriptions.OrderSubscription.SubscriptionId,
		}

		err := e.EsdbWorker.SubscribeAll(e.ctx, option)
		if err != nil {
			e.cancel()
			e.Log.Errorf("(esdbSubscribeAllWorker.SubscribeAll) err: {%v}", err)
		}
	}()
}
