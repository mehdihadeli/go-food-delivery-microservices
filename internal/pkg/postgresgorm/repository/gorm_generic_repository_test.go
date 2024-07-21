package repository

import (
	"context"
	"log"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/data"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	defaultLogger "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/defaultlogger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm"
	gorm2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/containers/testcontainer/gorm"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	_ "github.com/lib/pq" // postgres driver
)

// Product is a domain_events entity
type Product struct {
	ID          uuid.UUID
	Name        string
	Weight      int
	IsAvailable bool
}

// ProductGorm is DTO used to map Product entity to database
type ProductGorm struct {
	ID          uuid.UUID `gorm:"primaryKey;column:id"`
	Name        string    `gorm:"column:name"`
	Weight      int       `gorm:"column:weight"`
	IsAvailable bool      `gorm:"column:is_available"`
}

func (v *ProductGorm) TableName() string {
	return "products_gorm"
}

type gormGenericRepositoryTest struct {
	suite.Suite
	DB                             *gorm.DB
	productRepository              data.GenericRepository[*ProductGorm]
	productRepositoryWithDataModel data.GenericRepositoryWithDataModel[*ProductGorm, *Product]
	products                       []*ProductGorm
}

func TestGormGenericRepository(t *testing.T) {
	suite.Run(
		t,
		&gormGenericRepositoryTest{},
	)
}

func (c *gormGenericRepositoryTest) SetupSuite() {
	opts, err := gorm2.NewGormTestContainers(defaultLogger.GetLogger()).
		PopulateContainerOptions(context.Background(), c.T())
	c.Require().NoError(err)

	gormDB, err := postgresgorm.NewGorm(opts)
	c.Require().NoError(err)
	c.DB = gormDB

	err = migrationDatabase(gormDB)
	c.Require().NoError(err)

	c.productRepository = NewGenericGormRepository[*ProductGorm](gormDB)
	c.productRepositoryWithDataModel = NewGenericGormRepositoryWithDataModel[*ProductGorm, *Product](
		gormDB,
	)

	err = mapper.CreateMap[*ProductGorm, *Product]()
	if err != nil {
		log.Fatal(err)
	}

	err = mapper.CreateMap[*Product, *ProductGorm]()
	if err != nil {
		log.Fatal(err)
	}
}

func (c *gormGenericRepositoryTest) SetupTest() {
	p, err := seedData(context.Background(), c.DB)
	c.Require().NoError(err)
	c.products = p
}

func (c *gormGenericRepositoryTest) TearDownTest() {
	err := c.cleanupPostgresData()
	c.Require().NoError(err)
}

func (c *gormGenericRepositoryTest) Test_Add() {
	ctx := context.Background()

	product := &ProductGorm{
		ID:          uuid.NewV4(),
		Name:        gofakeit.Name(),
		Weight:      gofakeit.Number(100, 1000),
		IsAvailable: true,
	}

	err := c.productRepository.Add(ctx, product)
	c.Require().NoError(err)

	p, err := c.productRepository.GetById(ctx, product.ID)
	if err != nil {
		return
	}

	c.Assert().NotNil(p)
	c.Assert().Equal(product.ID, p.ID)
}

func (c *gormGenericRepositoryTest) Test_Add_With_Data_Model() {
	ctx := context.Background()

	product := &Product{
		ID:          uuid.NewV4(),
		Name:        gofakeit.Name(),
		Weight:      gofakeit.Number(100, 1000),
		IsAvailable: true,
	}

	err := c.productRepositoryWithDataModel.Add(ctx, product)
	c.Require().NoError(err)

	p, err := c.productRepositoryWithDataModel.GetById(ctx, product.ID)
	if err != nil {
		return
	}

	c.Assert().NotNil(p)
	c.Assert().Equal(product.ID, p.ID)
}

func (c *gormGenericRepositoryTest) Test_Get_By_Id() {
	ctx := context.Background()

	all, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	p := all.Items[0]

	testCases := []struct {
		Name         string
		ProductId    uuid.UUID
		ExpectResult *ProductGorm
	}{
		{
			Name:         p.Name,
			ProductId:    p.ID,
			ExpectResult: p,
		},
		{
			Name:         "NonExistingProduct",
			ProductId:    uuid.NewV4(),
			ExpectResult: nil,
		},
	}

	for _, s := range testCases {
		c.T().Run(s.Name, func(t *testing.T) {
			t.Parallel()
			res, err := c.productRepository.GetById(ctx, s.ProductId)
			if s.ExpectResult == nil {
				assert.Error(t, err)
				assert.True(t, customErrors.IsNotFoundError(err))
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, p.ID, res.ID)
			}
		})
	}
}

func (c *gormGenericRepositoryTest) Test_Get_By_Id_With_Data_Model() {
	ctx := context.Background()

	all, err := c.productRepositoryWithDataModel.GetAll(
		ctx,
		utils.NewListQuery(10, 1),
	)
	if err != nil {
		return
	}
	p := all.Items[0]

	testCases := []struct {
		Name         string
		ProductId    uuid.UUID
		ExpectResult *Product
	}{
		{
			Name:         p.Name,
			ProductId:    p.ID,
			ExpectResult: p,
		},
		{
			Name:         "NonExistingProduct",
			ProductId:    uuid.NewV4(),
			ExpectResult: nil,
		},
	}

	for _, s := range testCases {
		c.T().Run(s.Name, func(t *testing.T) {
			t.Parallel()
			res, err := c.productRepositoryWithDataModel.GetById(
				ctx,
				s.ProductId,
			)

			if s.ExpectResult == nil {
				assert.Error(t, err)
				assert.True(t, customErrors.IsNotFoundError(err))
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, p.ID, res.ID)
			}
		})
	}
}

func (c *gormGenericRepositoryTest) Test_Get_All() {
	ctx := context.Background()

	models, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	c.Assert().NotEmpty(models.Items)
}

func (c *gormGenericRepositoryTest) Test_Get_All_With_Data_Model() {
	ctx := context.Background()

	models, err := c.productRepositoryWithDataModel.GetAll(
		ctx,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models.Items)
}

func (c *gormGenericRepositoryTest) Test_Search() {
	ctx := context.Background()

	models, err := c.productRepository.Search(
		ctx,
		c.products[0].Name,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models.Items)
	c.Assert().Equal(len(models.Items), 1)
}

func (c *gormGenericRepositoryTest) Test_Search_With_Data_Model() {
	ctx := context.Background()

	models, err := c.productRepositoryWithDataModel.Search(
		ctx,
		c.products[0].Name,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models.Items)
	c.Assert().Equal(len(models.Items), 1)
}

func (c *gormGenericRepositoryTest) Test_Where() {
	ctx := context.Background()

	models, err := c.productRepository.GetByFilter(
		ctx,
		map[string]interface{}{"name": c.products[0].Name},
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models)
	c.Assert().Equal(len(models), 1)
}

func (c *gormGenericRepositoryTest) Test_Where_With_Data_Model() {
	ctx := context.Background()

	models, err := c.productRepositoryWithDataModel.GetByFilter(
		ctx,
		map[string]interface{}{"name": c.products[0].Name},
	)
	c.Require().NoError(err)

	c.Assert().NotEmpty(models)
	c.Assert().Equal(len(models), 1)
}

func (c *gormGenericRepositoryTest) Test_Update() {
	ctx := context.Background()

	products, err := c.productRepository.GetAll(ctx, utils.NewListQuery(10, 1))
	c.Require().NoError(err)

	product := products.Items[0]

	product.Name = "product2_updated"
	err = c.productRepository.Update(ctx, product)
	c.Require().NoError(err)

	single, err := c.productRepository.GetById(ctx, product.ID)
	c.Require().NoError(err)

	c.Assert().NotNil(single)
	c.Assert().Equal("product2_updated", single.Name)
}

func (c *gormGenericRepositoryTest) Test_Update_With_Data_Model() {
	ctx := context.Background()

	products, err := c.productRepositoryWithDataModel.GetAll(
		ctx,
		utils.NewListQuery(10, 1),
	)
	c.Require().NoError(err)

	product := products.Items[0]

	product.Name = "product2_updated"
	err = c.productRepositoryWithDataModel.Update(ctx, product)
	c.Require().NoError(err)

	single, err := c.productRepositoryWithDataModel.GetById(ctx, product.ID)
	c.Require().NoError(err)

	c.Assert().NotNil(single)
	c.Assert().Equal("product2_updated", single.Name)
}

//func Test_Delete(t *testing.T) {
//	ctx := context.Background()
//	repository, err := setupGenericGormRepository(ctx, t)
//
//	products, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
//	if err != nil {
//		t.Fatal(err)
//	}
//	product := products.Items[0]
//
//	err = repository.Delete(ctx, product.ID)
//	if err != nil {
//		return
//	}
//
//	single, err := repository.GetById(ctx, product.ID)
//	assert.Nil(t, single)
//}
//
//func Test_Delete_With_Data_Model(t *testing.T) {
//	ctx := context.Background()
//	repository, err := setupGenericGormRepositoryWithDataModel(ctx, t)
//
//	products, err := repository.GetAll(ctx, utils.NewListQuery(10, 1))
//	if err != nil {
//		t.Fatal(err)
//	}
//	product := products.Items[0]
//
//	err = repository.Delete(ctx, product.ID)
//	if err != nil {
//		return
//	}
//
//	single, err := repository.GetById(ctx, product.ID)
//	assert.Nil(t, single)
//}
//
//func Test_Count(t *testing.T) {
//	ctx := context.Background()
//	repository, err := setupGenericGormRepository(ctx, t)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	count := repository.Count(ctx)
//
//	assert.Equal(t, count, int64(2))
//}
//
//func Test_Count_With_Data_Model(t *testing.T) {
//	ctx := context.Background()
//	repository, err := setupGenericGormRepositoryWithDataModel(ctx, t)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	count := repository.Count(ctx)
//
//	assert.Equal(t, count, int64(2))
//}
//
//func Test_Find(t *testing.T) {
//	ctx := context.Background()
//	repository, err := setupGenericGormRepository(ctx, t)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	entities, err := repository.Find(
//		ctx,
//		specification.And(
//			specification.Equal("is_available", true),
//			specification.Equal("name", "seed_product1"),
//		),
//	)
//	if err != nil {
//		return
//	}
//	assert.Equal(t, len(entities), 1)
//}
//
//func Test_Find_With_Data_Model(t *testing.T) {
//	ctx := context.Background()
//	repository, err := setupGenericGormRepositoryWithDataModel(ctx, t)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	entities, err := repository.Find(
//		ctx,
//		specification.And(
//			specification.Equal("is_available", true),
//			specification.Equal("name", "seed_product1"),
//		),
//	)
//	if err != nil {
//		return
//	}
//	assert.Equal(t, len(entities), 1)
//}

func (c *gormGenericRepositoryTest) cleanupPostgresData() error {
	tables := []string{"products_gorm"}
	// Iterate over the tables and delete all records
	for _, table := range tables {
		err := c.DB.Exec("DELETE FROM " + table).Error

		return err
	}

	return nil
}

func migrationDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(ProductGorm{})
	if err != nil {
		return err
	}

	return nil
}

func seedData(ctx context.Context, db *gorm.DB) ([]*ProductGorm, error) {
	seedProducts := []*ProductGorm{
		{
			ID:          uuid.NewV4(),
			Name:        "seed_product1",
			Weight:      100,
			IsAvailable: true,
		},
		{
			ID:          uuid.NewV4(),
			Name:        "seed_product2",
			Weight:      100,
			IsAvailable: true,
		},
	}

	err := db.WithContext(ctx).Create(seedProducts).Error
	if err != nil {
		return nil, err
	}

	return seedProducts, nil
}
