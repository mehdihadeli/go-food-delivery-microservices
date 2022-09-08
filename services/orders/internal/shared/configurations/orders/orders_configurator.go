package orders

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations/order_module"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web"
	"github.com/pkg/errors"
	"net/http"
)

type OrdersServiceConfigurator interface {
	ConfigureProductsModule() error
}

type ordersServiceConfigurator struct {
	echoServer customEcho.EchoHttpServer
	grpcServer grpcServer.GrpcServer
	*infrastructure.InfrastructureConfiguration
}

func NewOrdersServiceConfigurator(infrastructureConfiguration *infrastructure.InfrastructureConfiguration, echoServer customEcho.EchoHttpServer, grpcServer grpcServer.GrpcServer) *ordersServiceConfigurator {
	return &ordersServiceConfigurator{InfrastructureConfiguration: infrastructureConfiguration, echoServer: echoServer, grpcServer: grpcServer}
}

func (c *ordersServiceConfigurator) ConfigureOrdersService(ctx context.Context) error {
	pc := order_module.NewOrdersModuleConfigurator(c.InfrastructureConfiguration, c.echoServer, c.grpcServer)
	err := pc.ConfigureOrdersModule(ctx)
	if err != nil {
		return errors.WithMessage(err, "[OrdersServiceConfigurator_ConfigureOrdersService.ConfigureOrdersModule] error in order module configurator")
	}

	err = c.migrateOrders()
	if err != nil {
		return errors.WithMessage(err, "[OrdersServiceConfigurator_ConfigureOrdersService.migrateOrders] error in the orders migration")
	}

	c.configSwagger()
	c.configMiddlewares(c.Metrics)

	c.echoServer.GetEchoInstance().GET("", func(ec echo.Context) error {
		return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", web.GetMicroserviceName(c.Cfg)))
	})

	return nil
}
