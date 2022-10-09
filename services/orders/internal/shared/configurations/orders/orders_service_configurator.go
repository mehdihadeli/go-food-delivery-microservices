package orders

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/labstack/echo/v4"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations"
	metrics2 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/orders/metrics"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/orders/rabbitmq"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/orders/subscription_all"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/contracts"
	"net/http"
)

type ordersServiceConfigurator struct {
	contracts.InfrastructureConfigurations
}

func NewOrdersServiceConfigurator(infrastructureConfiguration contracts.InfrastructureConfigurations) contracts.OrdersServiceConfigurator {
	return &ordersServiceConfigurator{InfrastructureConfigurations: infrastructureConfiguration}
}

func (c *ordersServiceConfigurator) ConfigureOrdersService(ctx context.Context) (contracts.OrderServiceConfigurations, error) {
	ordersServiceConfigurations := &ordersServiceConfigurations{}

	ordersServiceConfigurations.ordersGrpcServer = grpcServer.NewGrpcServer(c.Cfg().GRPC, c.Log(), c.Cfg().ServiceName, c.Metrics())
	ordersServiceConfigurations.ordersEchoServer = customEcho.NewEchoHttpServer(c.Cfg().Http, c.Log(), c.Cfg().ServiceName, c.Metrics())

	ordersServiceConfigurations.ordersEchoServer.RouteBuilder().RegisterRoutes(func(e *echo.Echo) {
		e.GET("", func(ec echo.Context) error {
			return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", c.Cfg().GetMicroserviceNameUpper()))
		})
	})

	// Orders Swagger Configs
	c.configSwagger(ordersServiceConfigurations.ordersEchoServer.RouteBuilder())

	// Orders Metrics Configs
	ordersMetrics, err := metrics2.ConfigOrdersMetrics(c.Cfg(), c.Metrics())
	if err != nil {
		return nil, err
	}
	ordersServiceConfigurations.ordersMetrics = ordersMetrics

	// Orders RabbitMQ Configs
	bus, err := rabbitmq.ConfigOrdersRabbitMQ(ctx, c.Cfg().RabbitMQ, c.InfrastructureConfigurations)
	if err != nil {
		return nil, err
	}
	ordersServiceConfigurations.ordersBus = bus

	// Orders SubscriptionsAll Configs
	esdbWorker, err := subscriptionAll.ConfigOrdersSubscriptionAllWorker(c.InfrastructureConfigurations, bus)
	if err != nil {
		return nil, err
	}
	ordersServiceConfigurations.esdbWorker = esdbWorker

	// Orders Product Module Configs
	pc := configurations.NewOrdersModuleConfigurator(c.InfrastructureConfigurations, ordersMetrics, bus, ordersServiceConfigurations.ordersEchoServer.RouteBuilder(), ordersServiceConfigurations.ordersGrpcServer.GrpcServiceBuilder())
	err = pc.ConfigureOrdersModule(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "[ordersServiceConfigurator_ConfigureOrdersService.NewOrdersModuleConfigurator] error in order module configurator")
	}

	return ordersServiceConfigurations, nil
}
