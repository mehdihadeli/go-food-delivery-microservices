package repositories

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	defaultLogger "github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/default_logger"
	gorm2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/containers/testcontainer/gorm"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/testfixture"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/data"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

//https://brunoscheufler.com/blog/2020-04-12-building-go-test-suites-using-testify

var (
	seedProductId1 = uuid.NewV4()
	seedProductId2 = uuid.NewV4()
)

// Define the custom testify suite
type ProductPostgresRepositoryTestSuite struct {
	suite.Suite
	productRepository data.ProductRepository
	ctx               context.Context
}

func TestProductPostgresRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &ProductPostgresRepositoryTestSuite{})
}

func (p *ProductPostgresRepositoryTestSuite) Test_Create_Product() {
	ctx := p.ctx

	product := &models.Product{
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
		ProductId:   uuid.NewV4(),
		Price:       gofakeit.Price(100, 1000),
		CreatedAt:   time.Now(),
	}

	createdProduct, err := p.productRepository.CreateProduct(ctx, product)
	require.NoError(p.T(), err)

	product, err = p.productRepository.GetProductById(ctx, createdProduct.ProductId)
	require.NoError(p.T(), err)

	assert.NotNil(p.T(), p)
	assert.Equal(p.T(), product.ProductId, createdProduct.ProductId)
}

func (p *ProductPostgresRepositoryTestSuite) Test_Update_Product() {
	ctx := p.ctx

	existingProduct, err := p.productRepository.GetProductById(ctx, seedProductId1)
	require.NoError(p.T(), err)
	require.NotNil(p.T(), existingProduct)

	existingProduct.Name = "test_update_product"
	_, err = p.productRepository.UpdateProduct(ctx, existingProduct)
	require.NoError(p.T(), err)

	updatedProduct, err := p.productRepository.GetProductById(ctx, seedProductId1)
	assert.Equal(p.T(), updatedProduct.Name, "test_update_product")
}

func (p *ProductPostgresRepositoryTestSuite) Test_Delete_Product() {
	ctx := p.ctx

	err := p.productRepository.DeleteProductByID(ctx, seedProductId1)
	require.NoError(p.T(), err)

	product, err := p.productRepository.GetProductById(ctx, seedProductId1)
	assert.NoError(p.T(), err)
	assert.Nil(p.T(), product)
}

func (p *ProductPostgresRepositoryTestSuite) Test_Get_Product() {
	ctx := p.ctx

	p.Run("Product Not Found", func() {
		// with subset test a new t will create for subset test
		res, err := p.productRepository.GetProductById(ctx, uuid.NewV4())
		require.NoError(p.T(), err)
		assert.Nil(p.T(), res)
	})

	p.Run("Get Product By ID", func() {
		res, err := p.productRepository.GetProductById(ctx, seedProductId1)
		require.NoError(p.T(), err)

		assert.NotNil(p.T(), res)
		assert.Equal(p.T(), res.ProductId, seedProductId1)
	})

	p.Run("Get All Products", func() {
		res, err := p.productRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
		require.NoError(p.T(), err)

		assert.Equal(p.T(), 2, len(res.Items))
	})
}

func (p *ProductPostgresRepositoryTestSuite) SetupSuite() {
	p.T().Log("SetupSuite")
}

func (p *ProductPostgresRepositoryTestSuite) SetupTest() {
	p.ctx = context.Background()
	p.T().Log("SetupTest")

	rep, err := setupTest(p.ctx, p)
	if err != nil {
		p.FailNowf("error in the setup repository", err.Error())
	}

	p.productRepository = rep
}

func (p *ProductPostgresRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	p.T().Log("BeforeTest")
}

func (p *ProductPostgresRepositoryTestSuite) AfterTest(suiteName, testName string) {
	p.T().Log("AfterTest")
}

func (p *ProductPostgresRepositoryTestSuite) TearDownSuite() {
	p.T().Log("TearDownSuite")
}

func (p *ProductPostgresRepositoryTestSuite) TearDownTest() {
	p.T().Log("TearDownTest")
	// cleanup test containers
	p.ctx.Done()
}

func setupTest(ctx context.Context, p *ProductPostgresRepositoryTestSuite) (data.ProductRepository, error) {
	gormDB, err := gorm2.NewGormTestContainers().Start(ctx, p.T())
	if err != nil {
		return nil, err
	}

	seedAndMigration(p, gormDB, []uuid.UUID{seedProductId1, seedProductId2})

	cfg, err := config.InitConfig(constants.Test)
	if err != nil {
		return nil, err
	}

	productRepository := NewPostgresProductRepository(defaultLogger.Logger, cfg, gormDB)
	return productRepository, nil
}

func seedAndMigration(p *ProductPostgresRepositoryTestSuite, gormDB *gorm.DB, productIds []uuid.UUID) {
	// migration
	err := gormDB.AutoMigrate(models.Product{})
	if err != nil {
		p.FailNowf("error in seed database", err.Error())
	}

	db, err := gormDB.DB()
	if err != nil {
		p.FailNowf("error in seed database", err.Error())
	}

	//https://github.com/go-testfixtures/testfixtures#templating
	// seed data
	var data []struct {
		Name        string
		ProductId   uuid.UUID
		Description string
	}

	for _, id := range productIds {
		data = append(data, struct {
			Name        string
			ProductId   uuid.UUID
			Description string
		}{
			Name:        gofakeit.Name(),
			Description: gofakeit.AdjectiveDescriptive(),
			ProductId:   id,
		})
	}

	err = testfixture.RunPostgresFixture(
		db,
		[]string{"db/fixtures/products"},
		map[string]interface{}{
			"Products": data,
		})
	if err != nil {
		p.FailNowf("error in seed database", err.Error())
	}
}
