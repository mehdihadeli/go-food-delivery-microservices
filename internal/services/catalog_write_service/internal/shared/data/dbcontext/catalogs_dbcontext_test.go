package dbcontext

import (
	"context"
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environment"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/external/fxlog"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/zap"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/migration/goose"
	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgresgorm"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/configurations/mappings"
	datamodel "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/data/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/models"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

// Define the suite
type DBContextTestSuite struct {
	suite.Suite
	items     []*datamodel.ProductDataModel
	dbContext *CatalogsGormDBContext
	app       *fxtest.App
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestDBContextTestSuite(t *testing.T) {
	suite.Run(t, new(DBContextTestSuite))
}

func (s *DBContextTestSuite) Test_FindProductByID() {
	s.Require().NotNil(s.dbContext)

	id := s.items[0].ProductId

	p, err := s.dbContext.FindProductByID(
		context.Background(),
		id,
	)
	s.Require().NoError(err)
	s.Require().NotNil(p)

	s.Assert().Equal(p.ProductId, id)
}

func (s *DBContextTestSuite) Test_DeleteProductByID() {
	s.Require().NotNil(s.dbContext)

	id := s.items[0].ProductId

	err := s.dbContext.DeleteProductByID(
		context.Background(),
		id,
	)
	s.Require().NoError(err)

	p, err := s.dbContext.FindProductByID(
		context.Background(),
		id,
	)

	s.Require().Nil(p)
	s.Require().Error(err)
}

func (s *DBContextTestSuite) Test_CreateProduct() {
	s.Require().NotNil(s.dbContext)

	item := &models.Product{
		ProductId:   uuid.NewV4(),
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(100, 1000),
	}

	res, err := s.dbContext.AddProduct(context.Background(), item)
	s.Require().NoError(err)

	p, err := s.dbContext.FindProductByID(
		context.Background(),
		item.ProductId,
	)
	s.Require().NoError(err)
	s.Require().NotNil(p)

	s.Assert().Equal(p.ProductId, item.ProductId)
	s.Assert().Equal(p.ProductId, res.ProductId)
}

func (s *DBContextTestSuite) Test_UpdateProduct() {
	s.Require().NotNil(s.dbContext)

	id := s.items[0].ProductId

	p, err := s.dbContext.FindProductByID(
		context.Background(),
		id,
	)
	s.Require().NoError(err)

	newName := gofakeit.Name()
	item := p
	item.Name = newName

	res, err := s.dbContext.UpdateProduct(context.Background(), item)
	s.Require().NoError(err)

	p2, err := s.dbContext.FindProductByID(
		context.Background(),
		id,
	)
	s.Require().NoError(err)

	s.Assert().Equal(item.Name, p2.Name)
	s.Assert().Equal(res.Name, p2.Name)
}

// TestSuite Hooks
func (s *DBContextTestSuite) SetupTest() {
	var gormDBContext *CatalogsGormDBContext

	app := fxtest.New(
		s.T(),
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		gormPostgres.Module,
		goose.Module,
		fx.Decorate(
			func(cfg *gormPostgres.GormOptions) (*gormPostgres.GormOptions, error) {
				// using sql-lite in-memory
				cfg.UseInMemory = true

				return cfg, nil
			},
		),
		fx.Provide(
			NewCatalogsDBContext,
		),
		fx.Populate(&gormDBContext),
	).RequireStart()

	s.dbContext = gormDBContext
	products := s.setupInMemoryContext()

	s.app = app
	s.items = products
}

func (s *DBContextTestSuite) SetupSuite() {
	err := mappings.ConfigureProductsMappings()
	s.Require().NoError(err)
}

func (s *DBContextTestSuite) setupInMemoryContext() []*datamodel.ProductDataModel {
	err := migrateGorm(s.dbContext.DB)
	s.Require().NoError(err)

	res, err := seedData(s.dbContext.DB)
	s.Require().NoError(errors.WrapIf(err, "error in seeding data in postgres"))

	return res
}

func migrateGorm(db *gorm.DB) error {
	err := db.AutoMigrate(&datamodel.ProductDataModel{})
	if err != nil {
		return err
	}

	return nil
}

func seedData(gormDB *gorm.DB) ([]*datamodel.ProductDataModel, error) {
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
	err := gormDB.CreateInBatches(products, len(products)).Error
	if err != nil {
		return nil, errors.Wrap(err, "error in seed database")
	}

	return products, nil
}
