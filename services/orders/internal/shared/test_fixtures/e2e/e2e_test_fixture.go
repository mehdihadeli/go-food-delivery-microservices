package e2e

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/store"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
)

type E2ETestFixture struct {
	Echo *echo.Echo
	*infrastructure.InfrastructureConfiguration
	V1                  *V1Groups
	GrpcServer          grpcServer.GrpcServer
	OrderAggregateStore store.AggregateStore[*aggregate.Order]
	Cleanup             func()
}

type V1Groups struct {
	OrdersGroup *echo.Group
}

func NewE2ETestFixture() *E2ETestFixture {
	cfg, _ := config.InitConfig("test")
	c := infrastructure.NewInfrastructureConfigurator(defaultLogger.Logger, cfg)
	infrastructures, _, cleanup := c.ConfigInfrastructures(context.Background())

	e := echo.New()

	v1Group := e.Group("/api/v1")
	ordersV1 := v1Group.Group("/orders")

	v1Groups := &V1Groups{OrdersGroup: ordersV1}

	eventStore := eventstroredb.NewEventStoreDbEventStore(infrastructures.Log, infrastructures.Esdb, infrastructures.EsdbSerializer)
	orderAggregateStore := eventstroredb.NewEventStoreAggregateStore[*aggregate.Order](infrastructures.Log, eventStore, infrastructures.EsdbSerializer)

	err := mediatr.ConfigOrdersMediator(orderAggregateStore, infrastructures)
	if err != nil {
		return nil
	}

	err = mappings.ConfigureMappings()
	if err != nil {
		return nil
	}

	grpcServer := grpcServer.NewGrpcServer(cfg.GRPC, defaultLogger.Logger)

	return &E2ETestFixture{
		Cleanup:                     cleanup,
		InfrastructureConfiguration: infrastructures,
		Echo:                        e,
		V1:                          v1Groups,
		OrderAggregateStore:         orderAggregateStore,
		GrpcServer:                  grpcServer,
	}
}
