package integration

import (
	"context"
	"testing"
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/bus"
	fxcontracts "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	gormPostgres "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/helpers/gormextensions"
	config2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/testfixture"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/config"
	datamodel "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/data/datamodels"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/app/test"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/data/dbcontext"
	productsService "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/grpc/genproto"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	rabbithole "github.com/michaelklishin/rabbit-hole"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

type IntegrationTestSharedFixture struct {
	Cfg                  *config.AppOptions
	Log                  logger.Logger
	Bus                  bus.Bus
	CatalogsDBContext    *dbcontext.CatalogsGormDBContext
	Container            fxcontracts.Container
	DbCleaner            dbcleaner.DbCleaner
	RabbitmqCleaner      *rabbithole.Client
	rabbitmqOptions      *config2.RabbitmqOptions
	Gorm                 *gorm.DB
	BaseAddress          string
	Items                []*datamodel.ProductDataModel
	ProductServiceClient productsService.ProductsServiceClient
}

func NewIntegrationTestSharedFixture(
	t *testing.T,
) *IntegrationTestSharedFixture {
	result := test.NewTestApp().Run(t)

	// https://github.com/michaelklishin/rabbit-hole
	rmqc, err := rabbithole.NewClient(
		result.RabbitmqOptions.RabbitmqHostOptions.HttpEndPoint(),
		result.RabbitmqOptions.RabbitmqHostOptions.UserName,
		result.RabbitmqOptions.RabbitmqHostOptions.Password)
	if err != nil {
		result.Logger.Error(
			errors.WrapIf(err, "error in creating rabbithole client"),
		)
	}

	shared := &IntegrationTestSharedFixture{
		Log:                  result.Logger,
		Container:            result.Container,
		Cfg:                  result.Cfg,
		RabbitmqCleaner:      rmqc,
		CatalogsDBContext:    result.CatalogsDBContext,
		Bus:                  result.Bus,
		rabbitmqOptions:      result.RabbitmqOptions,
		Gorm:                 result.Gorm,
		BaseAddress:          result.EchoHttpOptions.BasePathAddress(),
		ProductServiceClient: result.ProductServiceClient,
	}

	return shared
}

func (i *IntegrationTestSharedFixture) SetupTest() {
	i.Log.Info("SetupTest started")

	// migration will do in app configuration
	// seed data for our tests - app seed doesn't run in test environment
	res, err := seedDataManually(i.Gorm)
	if err != nil {
		i.Log.Error(errors.WrapIf(err, "error in seeding data in postgres"))
	}

	i.Items = res
}

func (i *IntegrationTestSharedFixture) TearDownTest() {
	i.Log.Info("TearDownTest started")

	// cleanup test containers with their hooks
	if err := i.cleanupRabbitmqData(); err != nil {
		i.Log.Error(errors.WrapIf(err, "error in cleanup rabbitmq data"))
	}

	if err := i.cleanupPostgresData(); err != nil {
		i.Log.Error(errors.WrapIf(err, "error in cleanup postgres data"))
	}
}

func (i *IntegrationTestSharedFixture) cleanupRabbitmqData() error {
	// https://github.com/michaelklishin/rabbit-hole
	// Get all queues
	queues, err := i.RabbitmqCleaner.ListQueuesIn(
		i.rabbitmqOptions.RabbitmqHostOptions.VirtualHost,
	)
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

func seedDataManually(gormDB *gorm.DB) ([]*datamodel.ProductDataModel, error) {
	products := []*datamodel.ProductDataModel{
		{
			Id:          uuid.NewV4(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
		{
			Id:          uuid.NewV4(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
	}

	err := gormDB.CreateInBatches(products, len(products)).Error
	if err != nil {
		return nil, errors.Wrap(err, "error in seed database")
	}

	return products, nil
}

func seedDataWithFixture(gormDB *gorm.DB) ([]*datamodel.ProductDataModel, error) {
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

	result, err := gormPostgres.Paginate[*datamodel.ProductDataModel, *datamodel.ProductDataModel](
		context.Background(),
		utils.NewListQuery(10, 1),
		gormDB,
	)

	return result.Items, nil
}
