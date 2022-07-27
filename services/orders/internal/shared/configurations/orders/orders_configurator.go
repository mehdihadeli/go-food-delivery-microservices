package orders

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web"
	"google.golang.org/grpc"
	"net/http"
)

type OrdersServiceConfigurator interface {
	ConfigureProductsModule() error
}

type ordersServiceConfigurator struct {
	log        logger.Logger
	cfg        *config.Config
	echo       *echo.Echo
	grpcServer *grpc.Server
}

func NewOrdersServiceConfigurator(log logger.Logger, cfg *config.Config, echo *echo.Echo, grpcServer *grpc.Server) *ordersServiceConfigurator {
	return &ordersServiceConfigurator{cfg: cfg, echo: echo, grpcServer: grpcServer, log: log}
}

func (c *ordersServiceConfigurator) ConfigureCatalogsService(ctx context.Context) (error, func()) {

	ic := infrastructure.NewInfrastructureConfigurator(c.log, c.cfg, c.echo, c.grpcServer)
	infrastructureConfigurations, err, infraCleanup := ic.ConfigInfrastructures(ctx)
	if err != nil {
		return err, nil
	}

	pc := configurations.NewOrdersModuleConfigurator(infrastructureConfigurations)
	err = pc.ConfigureOrdersModule(ctx)
	if err != nil {
		return err, nil
	}

	err = c.migrateOrders()
	if err != nil {
		return err, nil
	}

	c.echo.GET("", func(ec echo.Context) error {
		return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", web.GetMicroserviceName(c.cfg)))
	})

	return nil, infraCleanup
}
