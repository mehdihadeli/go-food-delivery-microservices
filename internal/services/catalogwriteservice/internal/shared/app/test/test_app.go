package test

import (
	"context"
	"os"
	"testing"
	"time"

	fxcontracts "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/grpc"
	config3 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/customecho/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	contracts2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/migration/contracts"
	gormPostgres "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/bus"
	config2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/containers/testcontainer/gorm"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/containers/testcontainer/rabbitmq"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/configurations/catalogs"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/data/dbcontext"
	productsService "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/grpc/genproto"

	"github.com/stretchr/testify/require"
	gorm2 "gorm.io/gorm"
)

type TestApp struct{}

type TestAppResult struct {
	Cfg                     *config.AppOptions
	Bus                     bus.RabbitmqBus
	Container               fxcontracts.Container
	Logger                  logger.Logger
	RabbitmqOptions         *config2.RabbitmqOptions
	EchoHttpOptions         *config3.EchoHttpOptions
	GormOptions             *gormPostgres.GormOptions
	Gorm                    *gorm2.DB
	ProductServiceClient    productsService.ProductsServiceClient
	GrpcClient              grpc.GrpcClient
	PostgresMigrationRunner contracts2.PostgresMigrationRunner
	CatalogsDBContext       *dbcontext.CatalogsGormDBContext
}

func NewTestApp() *TestApp {
	return &TestApp{}
}

func (a *TestApp) Run(t *testing.T) (result *TestAppResult) {
	lifetimeCtx := context.Background()

	// ref: https://github.com/uber-go/fx/blob/master/app_test.go
	appBuilder := NewCatalogsWriteTestApplicationBuilder(t)
	appBuilder.ProvideModule(catalogs.CatalogsServiceModule)

	appBuilder.Decorate(
		rabbitmq.RabbitmqContainerOptionsDecorator(t, lifetimeCtx),
	)
	appBuilder.Decorate(gorm.GormContainerOptionsDecorator(t, lifetimeCtx))

	testApp := appBuilder.Build()

	err := testApp.ConfigureCatalogs()
	if err != nil {
		testApp.Logger().Fatalf("Error in ConfigureCatalogs, %s", err)
	}

	err = testApp.MapCatalogsEndpoints()
	if err != nil {
		testApp.Logger().Fatalf("Error in MapCatalogsEndpoints, %s", err)
	}

	testApp.ResolveFunc(
		func(cfg *config.AppOptions,
			bus bus.RabbitmqBus,
			logger logger.Logger,
			rabbitmqOptions *config2.RabbitmqOptions,
			gormOptions *gormPostgres.GormOptions,
			gorm *gorm2.DB,
			catalogsDBContext *dbcontext.CatalogsGormDBContext,
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
				Gorm:                    gorm,
				CatalogsDBContext:       catalogsDBContext,
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

	err = testApp.Start(startCtx)
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
