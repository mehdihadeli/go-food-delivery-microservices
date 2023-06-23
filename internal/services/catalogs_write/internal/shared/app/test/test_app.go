package test

import (
	"context"
	"os"
	"testing"
	"time"

	gorm2 "gorm.io/gorm"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/metrics"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/bus"
	config2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/gorm"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/containers/testcontainer/rabbitmq"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/contracts/data"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/shared/configurations/catalogs"
)

type TestApp struct{}

type TestAppResult struct {
	Cfg                *config.AppOptions
	Bus                bus.RabbitmqBus
	Container          contracts.Container
	Logger             logger.Logger
	RabbitmqOptions    *config2.RabbitmqOptions
	GormOptions        *gormPostgres.GormOptions
	CatalogUnitOfWorks data.CatalogUnitOfWork
	ProductRepository  data.ProductRepository
	Gorm               *gorm2.DB
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
	appBuilder.Decorate(func(metrics *metrics.OtelMetrics) *metrics.OtelMetrics {
		metrics.Meter = nil
		return metrics
	})

	testApp := appBuilder.Build()

	testApp.ConfigureCatalogs()

	testApp.ResolveFunc(
		func(cfg *config.AppOptions,
			bus bus.RabbitmqBus,
			logger logger.Logger,
			rabbitmqOptions *config2.RabbitmqOptions,
			gormOptions *gormPostgres.GormOptions,
			catalogUnitOfWorks data.CatalogUnitOfWork,
			productRepository data.ProductRepository,
			gorm *gorm2.DB,
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
