package unittest

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/mocks"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	defaultLogger "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/defaultlogger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/external/gromlog"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/helpers/gormextensions"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/configurations/mappings"
	datamodel "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/data/datamodels"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/data/dbcontext"

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
	Bus              *mocks.Bus
	Tracer           trace.Tracer
	CatalogDBContext *dbcontext.CatalogsGormDBContext
	Ctx              context.Context
	dbFilePath       string
	dbFileName       string
}

func NewUnitTestSharedFixture(t *testing.T) *UnitTestSharedFixture {
	// we could use EmptyLogger if we don't want to log anything
	log := defaultLogger.GetLogger()
	cfg := &config.AppOptions{}

	// empty tracer, just for testing
	nopetracer := trace.NewNoopTracerProvider()
	testTracer := nopetracer.Tracer("test_tracer")

	unit := &UnitTestSharedFixture{
		Cfg:        cfg,
		Log:        log,
		Tracer:     testTracer,
		dbFileName: "sqlite.db",
	}

	return unit
}

func (c *UnitTestSharedFixture) BeginTx() {
	c.Log.Info("starting transaction")
	// seems when we `Begin` a transaction on gorm.DB (with SQLLite in-memory) our previous gormDB before transaction will remove and the new gormDB with tx will go on the memory
	tx := c.CatalogDBContext.DB().Begin()
	gormContext := gormextensions.SetTxToContext(c.Ctx, tx)
	c.Ctx = gormContext

	//// works on both transaction and none-transactional gormdbcontext
	//var productData []*datamodel.ProductDataModel
	//var productData2 []*datamodel.ProductDataModel
	//
	//s := c.CatalogDBContext.Find(&productData).Error
	//s2 := tx.Find(&productData2).Error
}

func (c *UnitTestSharedFixture) CommitTx() {
	tx := gormextensions.GetTxFromContextIfExists(c.Ctx)
	if tx != nil {
		c.Log.Info("committing transaction")
		tx.Commit()
	}
}

/// Shared Hooks

func (c *UnitTestSharedFixture) SetupSuite() {
	// this fix root working directory problem in our test environment inner our fixture
	environment.FixProjectRootWorkingDirectoryPath()
	projectRootDir := environment.GetProjectRootWorkingDirectory()

	c.dbFilePath = filepath.Join(projectRootDir, c.dbFileName)
}

func (c *UnitTestSharedFixture) TearDownSuite() {
}

func (c *UnitTestSharedFixture) SetupTest() {
	ctx := context.Background()
	c.Ctx = ctx

	c.setupBus()

	c.setupDB()

	err := mappings.ConfigureProductsMappings()
	c.Require().NoError(err)
}

func (c *UnitTestSharedFixture) TearDownTest() {
	err := c.cleanupDB()
	c.Require().NoError(err)

	mapper.ClearMappings()
}

func (c *UnitTestSharedFixture) setupBus() {
	// create new mocks
	bus := &mocks.Bus{}

	bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	c.Bus = bus
}

func (c *UnitTestSharedFixture) setupDB() {
	dbContext := c.createSQLLiteDBContext()
	c.CatalogDBContext = dbContext

	c.initDB(dbContext)
}

func (c *UnitTestSharedFixture) createSQLLiteDBContext() *dbcontext.CatalogsGormDBContext {
	// https://gorm.io/docs/connecting_to_the_database.html#SQLite
	// https://github.com/glebarez/sqlite
	// https://www.connectionstrings.com/sqlite/
	gormSQLLiteDB, err := gorm.Open(
		sqlite.Open(c.dbFilePath),
		&gorm.Config{
			Logger: gromlog.NewGormCustomLogger(defaultLogger.GetLogger()),
		})
	c.Require().NoError(err)

	dbContext := dbcontext.NewCatalogsDBContext(gormSQLLiteDB)

	return dbContext
}

func (c *UnitTestSharedFixture) initDB(dbContext *dbcontext.CatalogsGormDBContext) {
	// migrations for our database
	err := migrateGorm(dbContext)
	c.Require().NoError(err)

	// seed data for our tests
	items, err := seedDataManually(dbContext)
	c.Require().NoError(err)

	c.Products = items
}

func (c *UnitTestSharedFixture) cleanupDB() error {
	sqldb, _ := c.CatalogDBContext.DB().DB()
	e := sqldb.Close()
	c.Require().NoError(e)

	// removing sql-lite file
	err := os.Remove(c.dbFilePath)

	return err
}

func migrateGorm(dbContext *dbcontext.CatalogsGormDBContext) error {
	err := dbContext.DB().AutoMigrate(&datamodel.ProductDataModel{})
	if err != nil {
		return err
	}

	return nil
}

func seedDataManually(
	dbContext *dbcontext.CatalogsGormDBContext,
) ([]*datamodel.ProductDataModel, error) {
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

	// seed data
	err := dbContext.DB().CreateInBatches(products, len(products)).Error
	if err != nil {
		return nil, errors.Wrap(err, "error in seed database")
	}

	return products, nil
}
