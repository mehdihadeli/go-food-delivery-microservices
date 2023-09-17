package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/bus"
	config2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/testfixture"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/data"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/app/test"
	productsService "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/grpc/genproto"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

type IntegrationTestSharedFixture struct {
	Cfg                  *config.AppOptions
	Log                  logger.Logger
	Bus                  bus.Bus
	CatalogUnitOfWorks   data.CatalogUnitOfWork
	ProductRepository    data.ProductRepository
	Container            contracts.Container
	DbCleaner            dbcleaner.DbCleaner
	RabbitmqCleaner      *rabbithole.Client
	rabbitmqOptions      *config2.RabbitmqOptions
	Gorm                 *gorm.DB
	BaseAddress          string
	Items                []*models.Product
	ProductServiceClient productsService.ProductsServiceClient
}

func NewIntegrationTestSharedFixture(t *testing.T) *IntegrationTestSharedFixture {
	result := test.NewTestApp().Run(t)

	// https://github.com/michaelklishin/rabbit-hole
	rmqc, err := rabbithole.NewClient(
		fmt.Sprintf(result.RabbitmqOptions.RabbitmqHostOptions.HttpEndPoint()),
		result.RabbitmqOptions.RabbitmqHostOptions.UserName,
		result.RabbitmqOptions.RabbitmqHostOptions.Password)

	require.NoError(t, err)

	shared := &IntegrationTestSharedFixture{
		Log:                  result.Logger,
		Container:            result.Container,
		Cfg:                  result.Cfg,
		RabbitmqCleaner:      rmqc,
		ProductRepository:    result.ProductRepository,
		CatalogUnitOfWorks:   result.CatalogUnitOfWorks,
		Bus:                  result.Bus,
		rabbitmqOptions:      result.RabbitmqOptions,
		Gorm:                 result.Gorm,
		BaseAddress:          result.EchoHttpOptions.BasePathAddress(),
		ProductServiceClient: result.ProductServiceClient,
	}

	return shared
}

func (i *IntegrationTestSharedFixture) InitializeTest() {
	i.Log.Info("InitializeTest started")

	// seed data in each test
	res, err := seedData(i.Gorm)
	if err != nil {
		i.Log.Fatal(err)
	}
	i.Items = res
}

func (i *IntegrationTestSharedFixture) DisposeTest() {
	i.Log.Info("DisposeTest started")

	// cleanup test containers with their hooks
	if err := i.cleanupRabbitmqData(); err != nil {
		i.Log.Fatal(err)
	}

	if err := i.cleanupPostgresData(); err != nil {
		i.Log.Fatal(err)
	}
}

func (i *IntegrationTestSharedFixture) cleanupRabbitmqData() error {
	// https://github.com/michaelklishin/rabbit-hole
	// Get all queues
	queues, err := i.RabbitmqCleaner.ListQueuesIn(i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost)
	if err != nil {
		return err
	}

	// clear each queue
	for _, queue := range queues {
		_, err = i.RabbitmqCleaner.PurgeQueue(
			i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost,
			queue.Name,
		)
		return err
	}

	return nil
}

func (i *IntegrationTestSharedFixture) cleanupPostgresData() error {
	tables := []string{"products"}
	// Iterate over the tables and delete all records
	for _, table := range tables {
		err := i.Gorm.Exec("DELETE FROM " + table).Error
		return err
	}
	return nil
}

func seedData(gormDB *gorm.DB) ([]*models.Product, error) {
	products := []*models.Product{
		{
			ProductId:   uuid.NewV4(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
		{
			ProductId:   uuid.NewV4(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
	}

	// migration will do in app configuration
	// seed data
	err := gormDB.CreateInBatches(products, len(products)).Error
	if err != nil {
		return nil, errors.Wrap(err, "error in seed database")
	}
	return products, nil
}

func seedAndMigration(gormDB *gorm.DB) ([]*models.Product, error) {
	// migration
	err := gormDB.AutoMigrate(models.Product{})
	if err != nil {
		return nil, errors.WrapIf(err, "error in seed database")
	}

	db, err := gormDB.DB()
	if err != nil {
		return nil, errors.WrapIf(err, "error in seed database")
	}

	// https://github.com/go-testfixtures/testfixtures#templating
	// seed data
	var data []struct {
		Name        string
		ProductId   uuid.UUID
		Description string
	}

	f := []struct {
		Name        string
		ProductId   uuid.UUID
		Description string
	}{
		{gofakeit.Name(), uuid.NewV4(), gofakeit.AdjectiveDescriptive()},
		{gofakeit.Name(), uuid.NewV4(), gofakeit.AdjectiveDescriptive()},
	}

	data = append(data, f...)

	err = testfixture.RunPostgresFixture(
		db,
		[]string{"db/fixtures/products"},
		map[string]interface{}{
			"Products": data,
		})
	if err != nil {
		return nil, errors.WrapIf(err, "error in seed database")
	}

	result, err := gormPostgres.Paginate[*models.Product](
		context.Background(),
		utils.NewListQuery(10, 1),
		gormDB,
	)
	return result.Items, nil
}
