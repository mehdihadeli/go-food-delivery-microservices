package e2e

import (
	"context"
	"github.com/labstack/echo/v4"
	grpcServer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
)

type E2ETestFixture struct {
	Echo *echo.Echo
	*infrastructure.InfrastructureConfiguration
	V1                *V1Groups
	ProductRepository contracts.ProductRepository
	GrpcServer        grpcServer.GrpcServer
	Cleanup           func()
}

type V1Groups struct {
	ProductsGroup *echo.Group
}

func NewE2ETestFixture() *E2ETestFixture {
	cfg, _ := config.InitConfig("test")
	c := infrastructure.NewInfrastructureConfigurator(defaultLogger.Logger, cfg)
	infrastructures, _, cleanup := c.ConfigInfrastructures(context.Background())

	e := echo.New()

	v1Group := e.Group("/api/v1")
	productsV1 := v1Group.Group("/products")

	v1Groups := &V1Groups{ProductsGroup: productsV1}

	productRep := repositories.NewPostgresProductRepository(infrastructures.Log, cfg, infrastructures.Gorm.DB)

	err := mediatr.ConfigProductsMediator(productRep, infrastructures)
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
		ProductRepository:           productRep,
		GrpcServer:                  grpcServer,
	}
}
