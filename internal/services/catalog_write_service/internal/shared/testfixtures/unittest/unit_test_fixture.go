package unittest

import (
	"context"
	"testing"
	"time"

	mocks3 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/messaging/mocks"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/defaultlogger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/external/gromlog"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgresGorm/helpers"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/configurations/mappings"
	datamodel "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/data/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/data/dbcontext"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/glebarez/sqlite"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

type UnitTestSharedFixture struct {
	Cfg *config.AppOptions
	Log logger.Logger
	suite.Suite
	Products         []*datamodel.ProductDataModel
	Bus              *mocks3.Bus
	Tracer           trace.Tracer
	CatalogDBContext *dbcontext.CatalogsGormDBContext
	Ctx              context.Context
}

func NewUnitTestSharedFixture(t *testing.T) *UnitTestSharedFixture {
	// we could use EmptyLogger if we don't want to log anything
	log := defaultLogger.GetLogger()
	cfg := &config.AppOptions{}

	// empty tracer, just for testing
	nopetracer := trace.NewNoopTracerProvider()
	testTracer := nopetracer.Tracer("test_tracer")

	unit := &UnitTestSharedFixture{
		Cfg:    cfg,
		Log:    log,
		Tracer: testTracer,
	}

	return unit
}

// Shared Hooks

func (c *UnitTestSharedFixture) SetupTest() {
	// create new mocks
	bus := &mocks3.Bus{}

	bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	dbContext, products := c.createInMemoryDBContext()

	c.Bus = bus
	ctx := context.Background()
	c.Ctx = ctx
	c.CatalogDBContext = dbContext
	c.Products = products

	err := mappings.ConfigureProductsMappings()
	c.Require().NoError(err)
}

func (c *UnitTestSharedFixture) TearDownTest() {
	mapper.ClearMappings()
}

func (c *UnitTestSharedFixture) SetupSuite() {
}

func (c *UnitTestSharedFixture) TearDownSuite() {
}

func (c *UnitTestSharedFixture) createInMemoryDBContext() (*dbcontext.CatalogsGormDBContext, []*datamodel.ProductDataModel) {
	inMemoryGormDB, err := createInMemoryDB()
	c.Require().NoError(err)

	dbContext := dbcontext.NewCatalogsDBContext(inMemoryGormDB, c.Log)

	err = migrateGorm(dbContext)
	c.Require().NoError(err)

	items, err := seedData(dbContext)
	c.Require().NoError(err)

	return dbContext, items
}

func (c *UnitTestSharedFixture) BeginTx() {
	c.Log.Info("starting transaction")
	tx := c.CatalogDBContext.Begin()
	gormContext := helpers.SetTxToContext(c.Ctx, tx)
	c.Ctx = gormContext
}

func (c *UnitTestSharedFixture) CommitTx() {
	tx, err := helpers.GetTxFromContext(c.Ctx)
	c.Require().NoError(err)
	c.Log.Info("committing transaction")
	tx.Commit()
}

func createInMemoryDB() (*gorm.DB, error) {
	// https://gorm.io/docs/connecting_to_the_database.html#SQLite
	// https://github.com/glebarez/sqlite
	// https://www.connectionstrings.com/sqlite/
	db, err := gorm.Open(
		sqlite.Open(":memory:"),
		&gorm.Config{
			Logger: gromlog.NewGormCustomLogger(defaultLogger.GetLogger()),
		})

	return db, err
}

func migrateGorm(dbContext *dbcontext.CatalogsGormDBContext) error {
	err := dbContext.DB.AutoMigrate(&datamodel.ProductDataModel{})
	if err != nil {
		return err
	}

	return nil
}

func seedData(
	dbContext *dbcontext.CatalogsGormDBContext,
) ([]*datamodel.ProductDataModel, error) {
	products := []*datamodel.ProductDataModel{
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
	err := dbContext.DB.CreateInBatches(products, len(products)).Error
	if err != nil {
		return nil, errors.Wrap(err, "error in seed database")
	}

	return products, nil
}
