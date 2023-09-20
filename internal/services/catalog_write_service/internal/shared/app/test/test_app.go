package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	config3 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	contracts2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/migration/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/migration/goose"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/bus"
	config2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/gorm"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/rabbitmq"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/data"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/configurations/catalogs"
	productsService "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/grpc/genproto"

	"github.com/stretchr/testify/require"
	gorm2 "gorm.io/gorm"
)

type TestApp struct{}

type TestAppResult struct {
	Cfg                     *config.AppOptions
	Bus                     bus.RabbitmqBus
	Container               contracts.Container
	Logger                  logger.Logger
	RabbitmqOptions         *config2.RabbitmqOptions
	EchoHttpOptions         *config3.EchoHttpOptions
	GormOptions             *gormPostgres.GormOptions
	CatalogUnitOfWorks      data.CatalogUnitOfWork
	ProductRepository       data.ProductRepository
	Gorm                    *gorm2.DB
	ProductServiceClient    productsService.ProductsServiceClient
	GrpcClient              grpc.GrpcClient
	PostgresMigrationRunner contracts2.PostgresMigrationRunner
}

func NewTestApp() *TestApp {
	return &TestApp{}
}

func (a *TestApp) Run(t *testing.T) (result *TestAppResult) {
	lifetimeCtx := context.Background()

	// ref: https://github.com/uber-go/fx/blob/master/app_test.go
	appBuilder := NewCatalogsWriteTestApplicationBuilder(t)
	appBuilder.ProvideModule(catalogs.CatalogsServiceModule)
	appBuilder.ProvideModule(goose.Module)

	appBuilder.Decorate(rabbitmq.RabbitmqContainerOptionsDecorator(t, lifetimeCtx))
	appBuilder.Decorate(gorm.GormContainerOptionsDecorator(t, lifetimeCtx))

	testApp := appBuilder.Build()

	testApp.ConfigureCatalogs()

	testApp.MapCatalogsEndpoints()

	testApp.ResolveFunc(
		func(cfg *config.AppOptions,
			bus bus.RabbitmqBus,
			logger logger.Logger,
			rabbitmqOptions *config2.RabbitmqOptions,
			gormOptions *gormPostgres.GormOptions,
			catalogUnitOfWorks data.CatalogUnitOfWork,
			productRepository data.ProductRepository,
			gorm *gorm2.DB,
			echoOptions *config3.EchoHttpOptions,
			grpcClient grpc.GrpcClient,
			postgresMigrationRunner contracts2.PostgresMigrationRunner,
		) {
			grpcConnection := grpcClient.GetGrpcConnection()

			result = &TestAppResult{
				Bus:                     bus,
				Cfg:                     cfg,
				Container:               testApp,
				Logger:                  logger,
				RabbitmqOptions:         rabbitmqOptions,
				GormOptions:             gormOptions,
				ProductRepository:       productRepository,
				CatalogUnitOfWorks:      catalogUnitOfWorks,
				Gorm:                    gorm,
				EchoHttpOptions:         echoOptions,
				PostgresMigrationRunner: postgresMigrationRunner,
				ProductServiceClient: productsService.NewProductsServiceClient(
					grpcConnection,
				),
				GrpcClient: grpcClient,
			}
		},
	)
	// we need a longer timout for up and running our testcontainers
	duration := time.Second * 300

	// short timeout for handling start hooks and setup dependencies
	startCtx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	err := testApp.Start(startCtx)
	if err != nil {
		t.Errorf("Error starting, err: %v", err)
		os.Exit(1)
	}

	// waiting for grpc endpoint becomes ready in the given timeout
	err = result.GrpcClient.WaitForAvailableConnection()
	require.NoError(t, err)

	t.Cleanup(func() {
		// short timeout for handling stop hooks
		stopCtx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		err = testApp.Stop(stopCtx)
		require.NoError(t, err)
	})

	return
}
