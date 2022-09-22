package e2e

import (
	"context"
	"github.com/labstack/echo/v4"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	bus2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/bus"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/configurations/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/configurations/infrastructure"
	"net/http/httptest"
)

type E2ETestFixture struct {
	Echo *echo.Echo
	*infrastructure.InfrastructureConfigurations
	V1          *V1Groups
	GrpcServer  grpcServer.GrpcServer
	HttpServer  *httptest.Server
	RabbitMQBus bus2.Bus
	ctx         context.Context
	cancel      context.CancelFunc
	Cleanup     func()
}

type V1Groups struct {
	ProductsGroup *echo.Group
}

func NewE2ETestFixture() *E2ETestFixture {
	ctx, cancel := context.WithCancel(context.Background())
	cfg, _ := config.InitConfig("test")
	c := infrastructure.NewInfrastructureConfigurator(defaultLogger.Logger, cfg)
	infrastructures, _, cleanup := c.ConfigInfrastructures(context.Background())
	echo := echo.New()

	v1Group := echo.Group("/api/v1")
	productsV1 := v1Group.Group("/products")

	v1Groups := &V1Groups{ProductsGroup: productsV1}

	err := mediatr.ConfigProductsMediator(infrastructures)
	if err != nil {
		cancel()
		return nil
	}

	err = mappings.ConfigureMappings()
	if err != nil {
		cancel()
		return nil
	}

	grpcServer := grpcServer.NewGrpcServer(cfg.GRPC, defaultLogger.Logger)
	httpServer := httptest.NewServer(echo)

	rabbitmqBus := bus.NewRabbitMQBus(infrastructures.Log, infrastructures.Consumers)

	return &E2ETestFixture{
		Cleanup: func() {
			cancel()
			cleanup()
			grpcServer.GracefulShutdown()
			echo.Shutdown(ctx)
			httpServer.Close()
			rabbitmqBus.Stop(ctx)
		},
		InfrastructureConfigurations: infrastructures,
		Echo:                         echo,
		V1:                           v1Groups,
		GrpcServer:                   grpcServer,
		HttpServer:                   httpServer,
		RabbitMQBus:                  rabbitmqBus,
		ctx:                          ctx,
		cancel:                       cancel,
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
		err := e.RabbitMQBus.Start(e.ctx)
		if err != nil {
			e.cancel()
			e.Log.Errorf("(RabbitMQBus.Start) err: {%v}", err)
		}
	}()
}
