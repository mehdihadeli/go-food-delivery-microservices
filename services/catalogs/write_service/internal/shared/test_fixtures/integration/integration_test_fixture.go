package integration

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/configurations/mappings"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/data/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations/infrastructure"
)

type IntegrationTestFixture struct {
	*infrastructure.InfrastructureConfiguration
	ProductRepository contracts.ProductRepository
	ctx               context.Context
	cancel            context.CancelFunc
	Cleanup           func()
}

func NewIntegrationTestFixture() *IntegrationTestFixture {
	cfg, _ := config.InitConfig(constants.Test)

	ctx, cancel := context.WithCancel(context.Background())
	c := infrastructure.NewInfrastructureConfigurator(defaultLogger.Logger, cfg)
	infrastructures, _, cleanup := c.ConfigInfrastructures(context.Background())

	productRep := repositories.NewPostgresProductRepository(infrastructures.Log, cfg, infrastructures.Gorm.DB)

	err := mappings.ConfigureMappings()
	if err != nil {
		cancel()
		return nil
	}

	return &IntegrationTestFixture{
		Cleanup: func() {
			cancel()
			cleanup()
		},
		InfrastructureConfiguration: infrastructures,
		ProductRepository:           productRep,
		ctx:                         ctx,
		cancel:                      cancel,
	}
}

func (e *IntegrationTestFixture) Run() {

}
