package test

import (
	"context"
	"os"
	"testing"
	"time"

	gorm2 "gorm.io/gorm"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/data"
	productsService "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/grpc/proto/service_clients"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/configurations/catalogs"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc"
	config3 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/bus"
	config2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/gorm"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/rabbitmq"
)

type TestApp struct{}

type TestAppResult struct {
	Cfg                  *config.AppOptions
	Bus                  bus.RabbitmqBus
	Container            contracts.Container
	Logger               logger.Logger
	RabbitmqOptions      *config2.RabbitmqOptions
	EchoHttpOptions      *config3.EchoHttpOptions
	GormOptions          *gormPostgres.GormOptions
	CatalogUnitOfWorks   data.CatalogUnitOfWork
	ProductRepository    data.ProductRepository
	Gorm                 *gorm2.DB
	ProductServiceClient productsService.ProductsServiceClient
}

func NewTestApp() *TestApp {
	return &TestApp{}
}

func (a *TestApp) Run(t *testing.T) (result *TestAppResult) {
	lifetimeCtx := context.Background()

	// ref: https://github.com/uber-go/fx/blob/master/app_test.go
	appBuilder := NewCatalogsWriteTestApplicationBuilder(t)
	appBuilder.ProvideModule(catalogs.CatalogsServiceModule)
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
		) {
			result = &TestAppResult{
				Bus:                bus,
				Cfg:                cfg,
				Container:          testApp,
				Logger:             logger,
				RabbitmqOptions:    rabbitmqOptions,
				GormOptions:        gormOptions,
				ProductRepository:  productRepository,
				CatalogUnitOfWorks: catalogUnitOfWorks,
				Gorm:               gorm,
				EchoHttpOptions:    echoOptions,
				ProductServiceClient: productsService.NewProductsServiceClient(
					grpcClient.GetGrpcConnection(),
				),
			}
		},
	)
	duration := time.Second * 20

	// short timeout for handling start hooks and setup dependencies
	startCtx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	err := testApp.Start(startCtx)
	if err != nil {
		os.Exit(1)
	}

	t.Cleanup(func() {
		// short timeout for handling stop hooks
		stopCtx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		_ = testApp.Stop(stopCtx)
	})

	return
}
