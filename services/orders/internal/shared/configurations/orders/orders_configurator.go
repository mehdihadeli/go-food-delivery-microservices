package orders

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	customEcho "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/custom_echo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/configurations"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/configurations/infrastructure"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/web"
	"net/http"
)

type OrdersServiceConfigurator interface {
	ConfigureProductsModule() error
}

type ordersServiceConfigurator struct {
	log        logger.Logger
	cfg        *config.Config
	ehcoServer customEcho.EchoHttpServer
	grpcServer grpcServer.GrpcServer
}

func NewOrdersServiceConfigurator(log logger.Logger, cfg *config.Config, ehcoServer customEcho.EchoHttpServer, grpcServer grpcServer.GrpcServer) *ordersServiceConfigurator {
	return &ordersServiceConfigurator{cfg: cfg, ehcoServer: ehcoServer, grpcServer: grpcServer, log: log}
}

func (c *ordersServiceConfigurator) ConfigureCatalogsService(ctx context.Context) (error, func()) {

	ic := infrastructure.NewInfrastructureConfigurator(c.log, c.cfg, c.ehcoServer, c.grpcServer)
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

	c.ehcoServer.GetEchoInstance().GET("", func(ec echo.Context) error {
		return ec.String(http.StatusOK, fmt.Sprintf("%s is running...", web.GetMicroserviceName(c.cfg)))
	})

	return nil, infraCleanup
}
