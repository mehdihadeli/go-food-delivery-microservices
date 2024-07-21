//go:build unit
// +build unit

package gormdbcontext

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/external/fxlog"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/zap"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mapper"
	gormPostgres "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/contracts"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/scopes"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/goccy/go-json"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

// ProductDataModel data model
type ProductDataModel struct {
	Id          uuid.UUID `gorm:"primaryKey"`
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time `gorm:"default:current_timestamp"`
	UpdatedAt   time.Time
	// for soft delete - https://gorm.io/docs/delete.html#Soft-Delete
	gorm.DeletedAt
}

// TableName overrides the table name used by ProductDataModel to `products` - https://gorm.io/docs/conventions.html#TableName
func (p *ProductDataModel) TableName() string {
	return "products"
}

func (p *ProductDataModel) String() string {
	j, _ := json.Marshal(p)

	return string(j)
}

// Product model
type Product struct {
	Id          uuid.UUID
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Define the suite
type GormDBContextTestSuite struct {
	suite.Suite
	items      []*ProductDataModel
	dbContext  contracts.GormDBContext
	app        *fxtest.App
	dbFilePath string
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestGormDBContext(t *testing.T) {
	suite.Run(t, new(GormDBContextTestSuite))
}

func (s *GormDBContextTestSuite) Test_FindProductByID() {
	s.Require().NotNil(s.dbContext)

	id := s.items[0].Id

	p, err := FindModelByID[*ProductDataModel, *Product](
		context.Background(),
		s.dbContext,
		id,
	)
	s.Require().NoError(err)
	s.Require().NotNil(p)

	s.Assert().Equal(p.Id, id)
}

func (s *GormDBContextTestSuite) Test_ExistsProductByID() {
	s.Require().NotNil(s.dbContext)

	id := s.items[0].Id

	exist := Exists[*ProductDataModel](
		context.Background(),
		s.dbContext,
		id,
	)

	s.Require().True(exist)
}

func (s *GormDBContextTestSuite) Test_NoneExistsProductByID() {
	s.Require().NotNil(s.dbContext)

	id := uuid.NewV4()

	exist := Exists[*ProductDataModel](
		context.Background(),
		s.dbContext,
		id,
	)

	s.Require().False(exist)
}

func (s *GormDBContextTestSuite) Test_DeleteProductByID() {
	s.Require().NotNil(s.dbContext)

	id := s.items[0].Id

	err := DeleteDataModelByID[*ProductDataModel](
		context.Background(),
		s.dbContext,
		id,
	)
	s.Require().NoError(err)

	p, err := FindModelByID[*ProductDataModel, *Product](
		context.Background(),
		s.dbContext,
		id,
	)
	s.Require().Error(err)
	s.Require().Nil(p)

	// https://gorm.io/docs/delete.html#Find-soft-deleted-records
	var softDeletedProduct *ProductDataModel
	s.dbContext.DB().Scopes(scopes.FilterAllItemsWithSoftDeleted).First(&softDeletedProduct, id)
	s.Require().NotNil(softDeletedProduct)

	var deletedCount int64
	var allCount int64

	// https://gorm.io/docs/advanced_query.html#Count
	s.dbContext.DB().Model(&ProductDataModel{}).Scopes(scopes.FilterAllItemsWithSoftDeleted).Count(&allCount)
	s.Equal(allCount, int64(2))

	s.dbContext.DB().Model(&ProductDataModel{}).Scopes(scopes.SoftDeleted).Count(&deletedCount)
	s.Equal(deletedCount, int64(1))
}

func (s *GormDBContextTestSuite) Test_CreateProduct() {
	s.Require().NotNil(s.dbContext)

	item := &Product{
		Id:          uuid.NewV4(),
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(100, 1000),
	}

	res, err := AddModel[*ProductDataModel, *Product](context.Background(), s.dbContext, item)
	s.Require().NoError(err)

	p, err := FindModelByID[*ProductDataModel, *Product](
		context.Background(),
		s.dbContext,
		item.Id,
	)
	s.Require().NoError(err)
	s.Require().NotNil(p)

	s.Assert().Equal(p.Id, item.Id)
	s.Assert().Equal(p.Id, res.Id)
}

func (s *GormDBContextTestSuite) Test_UpdateProduct() {
	s.Require().NotNil(s.dbContext)

	id := s.items[0].Id

	p, err := FindModelByID[*ProductDataModel, *Product](
		context.Background(),
		s.dbContext,
		id,
	)
	s.Require().NoError(err)

	newName := gofakeit.Name()
	item := p
	item.Name = newName

	res, err := UpdateModel[*ProductDataModel, *Product](context.Background(), s.dbContext, item)
	s.Require().NoError(err)

	p2, err := FindModelByID[*ProductDataModel, *Product](
		context.Background(),
		s.dbContext,
		id,
	)
	s.Require().NoError(err)

	s.Assert().Equal(item.Name, p2.Name)
	s.Assert().Equal(res.Name, p2.Name)
}

// TestSuite Hooks

func (s *GormDBContextTestSuite) SetupTest() {
	err := ConfigureProductsMappings()
	s.Require().NoError(err)

	var gormDBContext contracts.GormDBContext
	var gormOptions *gormPostgres.GormOptions

	app := fxtest.New(
		s.T(),
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		gormPostgres.Module,
		fx.Provide(NewGormDBContext),
		fx.Decorate(
			func(cfg *gormPostgres.GormOptions) (*gormPostgres.GormOptions, error) {
				// using sql-lite with a database file
				cfg.UseSQLLite = true

				return cfg, nil
			},
		),
		fx.Populate(&gormDBContext),
		fx.Populate(&gormOptions),
	).RequireStart()

	s.app = app
	s.dbContext = gormDBContext
	s.dbFilePath = gormOptions.Dns()

	s.initDB()
}

func (s *GormDBContextTestSuite) TearDownTest() {
	err := s.cleanupDB()
	s.Require().NoError(err)

	mapper.ClearMappings()

	s.app.RequireStop()
}

func (s *GormDBContextTestSuite) initDB() {
	err := migrateGorm(s.dbContext.DB())
	s.Require().NoError(err)

	products, err := seedData(s.dbContext.DB())
	s.Require().NoError(err)

	s.items = products
}

func (s *GormDBContextTestSuite) cleanupDB() error {
	sqldb, _ := s.dbContext.DB().DB()
	e := sqldb.Close()
	s.Require().NoError(e)

	// removing sql-lite file
	err := os.Remove(s.dbFilePath)

	return err
}

func migrateGorm(db *gorm.DB) error {
	err := db.AutoMigrate(&ProductDataModel{})
	if err != nil {
		return err
	}

	return nil
}

func seedData(gormDB *gorm.DB) ([]*ProductDataModel, error) {
	products := []*ProductDataModel{
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
	err := gormDB.CreateInBatches(products, len(products)).Error
	if err != nil {
		return nil, errors.Wrap(err, "error in seed database")
	}

	return products, nil
}

func ConfigureProductsMappings() error {
	err := mapper.CreateMap[*ProductDataModel, *Product]()
	if err != nil {
		return err
	}

	err = mapper.CreateMap[*Product, *ProductDataModel]()
	if err != nil {
		return err
	}

	return nil
}
