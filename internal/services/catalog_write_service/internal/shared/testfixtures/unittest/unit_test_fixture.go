package unittest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environment"
	mocks3 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/messaging/mocks"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/defaultlogger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/external/gromlog"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgresgorm/helpers"
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
	projectRootDir   string
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

// Shared Hooks

func (c *UnitTestSharedFixture) SetupTest() {
	ctx := context.Background()
	c.Ctx = ctx

	// create new mocks
	bus := &mocks3.Bus{}

	bus.On("PublishMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	dbContext, products := c.createSQLLiteDBContext()

	c.Bus = bus
	c.CatalogDBContext = dbContext
	c.Products = products

	err := mappings.ConfigureProductsMappings()
	c.Require().NoError(err)
}

func (c *UnitTestSharedFixture) TearDownTest() {
	err := c.dropSQLLiteDB()
	c.Require().NoError(err)
	mapper.ClearMappings()
}

func (c *UnitTestSharedFixture) SetupSuite() {
	environment.FixProjectRootWorkingDirectoryPath()
	c.projectRootDir = environment.GetProjectRootWorkingDirectory()
}

func (c *UnitTestSharedFixture) TearDownSuite() {
}

func (c *UnitTestSharedFixture) createSQLLiteDBContext() (*dbcontext.CatalogsGormDBContext, []*datamodel.ProductDataModel) {
	sqlLiteGormDB, err := c.createSQLLiteDB()
	c.Require().NoError(err)

	dbContext := dbcontext.NewCatalogsDBContext(sqlLiteGormDB, c.Log)

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

	var productData []*datamodel.ProductDataModel
	var productData2 []*datamodel.ProductDataModel

	s := c.CatalogDBContext.Find(&productData).Error
	s2 := tx.Find(&productData2).Error
	fmt.Println(s)
	fmt.Println(s2)
}

func (c *UnitTestSharedFixture) CommitTx() {
	tx := helpers.GetTxFromContextIfExists(c.Ctx)
	if tx != nil {
		c.Log.Info("committing transaction")
		tx.Commit()
	}
}

func (c *UnitTestSharedFixture) createSQLLiteDB() (*gorm.DB, error) {
	dbFilePath := filepath.Join(c.projectRootDir, c.dbFileName)

	// https://gorm.io/docs/connecting_to_the_database.html#SQLite
	// https://github.com/glebarez/sqlite
	// https://www.connectionstrings.com/sqlite/
	gormDB, err := gorm.Open(
		sqlite.Open(dbFilePath),
		&gorm.Config{
			Logger: gromlog.NewGormCustomLogger(defaultLogger.GetLogger()),
		})

	return gormDB, err
}

func (c *UnitTestSharedFixture) dropSQLLiteDB() error {
	sqldb, _ := c.CatalogDBContext.DB.DB()
	e := sqldb.Close()
	c.Require().NoError(e)

	dbFilePath := filepath.Join(c.projectRootDir, c.dbFileName)
	err := os.Remove(dbFilePath)

	return err
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
