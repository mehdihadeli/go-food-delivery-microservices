package repositories

import (
	"context"
	"testing"
	"time"

	gormPostgres "github.com/mehdihadeli/store-golang-microservice-sample/pkg/gorm_postgres"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/testfixture"

	"github.com/brianvoe/gofakeit/v6"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	defaultLogger "github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/default_logger"
	gorm2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/containers/testcontainer/gorm"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/data"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/models"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// https://brunoscheufler.com/blog/2020-04-12-building-go-test-suites-using-testify

var items []*models.Product

// Define the custom testify suite
type productPostgresRepositoryTestSuite struct {
	suite.Suite
	productRepository data.ProductRepository
	ctx               context.Context
}

func TestProductPostgresRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &productPostgresRepositoryTestSuite{})
}

func (p *productPostgresRepositoryTestSuite) Test_Create_Product_Should_Create_New_Product_In_DB() {
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
	p.NoError(err)

	p.NotNil(p)
	p.Equal(product.ProductId, createdProduct.ProductId)
}

func (p *productPostgresRepositoryTestSuite) Test_Update_Product_Should_Update_Existing_Product_In_DB() {
	ctx := p.ctx

	id := items[0].ProductId
	existingProduct, err := p.productRepository.GetProductById(ctx, id)
	p.Require().NoError(err)
	p.Require().NotNil(existingProduct)

	existingProduct.Name = "test_update_product"
	_, err = p.productRepository.UpdateProduct(ctx, existingProduct)
	p.Require().NoError(err)

	updatedProduct, err := p.productRepository.GetProductById(ctx, id)
	p.Equal(updatedProduct.Name, "test_update_product")
}

func (p *productPostgresRepositoryTestSuite) Test_Delete_Product_Should_Delete_Existing_Product_In_DB() {
	ctx := p.ctx

	id := items[0].ProductId

	err := p.productRepository.DeleteProductByID(ctx, id)
	p.Require().NoError(err)

	product, err := p.productRepository.GetProductById(ctx, id)

	p.Error(err)
	p.True(customErrors.IsNotFoundError(err))
	p.Nil(product)
}

func (p *productPostgresRepositoryTestSuite) Test_Get_Product() {
	ctx := p.ctx
	id := items[0].ProductId

	p.Run("Should_Return_NotFound_Error_When_Item_DoesNot_Exists", func() {
		// with subset test a new t will create for subset test
		res, err := p.productRepository.GetProductById(ctx, uuid.NewV4())

		p.Error(err)
		p.True(customErrors.IsNotFoundError(err))
		p.Nil(res)
	})

	p.Run("Should_Get_Existing_Product_From_DB", func() {
		res, err := p.productRepository.GetProductById(ctx, id)
		p.Require().NoError(err)

		p.NotNil(res)
		p.Equal(res.ProductId, id)
	})

	p.Run("Should_Get_All_Existing_Products_From_DB", func() {
		res, err := p.productRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
		p.Require().NoError(err)

		p.Equal(2, len(res.Items))
	})
}

func (p *productPostgresRepositoryTestSuite) SetupSuite() {
	p.T().Log("SetupSuite")
}

func (p *productPostgresRepositoryTestSuite) SetupTest() {
	p.ctx = context.Background()
	p.T().Log("SetupTest")

	rep, err := setupTest(p.ctx, p)
	if err != nil {
		p.FailNowf("error in the setup repository", err.Error())
	}

	p.productRepository = rep
}

func (p *productPostgresRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	p.T().Log("BeforeTest")
}

func (p *productPostgresRepositoryTestSuite) AfterTest(suiteName, testName string) {
	p.T().Log("AfterTest")
}

func (p *productPostgresRepositoryTestSuite) TearDownSuite() {
	p.T().Log("TearDownSuite")
}

func (p *productPostgresRepositoryTestSuite) TearDownTest() {
	p.T().Log("TearDownTest")
	// cleanup test containers
	p.ctx.Done()
}

func setupTest(ctx context.Context, p *productPostgresRepositoryTestSuite) (data.ProductRepository, error) {
	gormDB, err := gorm2.NewGormTestContainers().Start(ctx, p.T())
	if err != nil {
		return nil, err
	}

	seedAndMigration(p, gormDB)

	productRepository := NewPostgresProductRepository(defaultLogger.Logger, gormDB)

	return productRepository, nil
}

func seedAndMigration(p *productPostgresRepositoryTestSuite, gormDB *gorm.DB) {
	// migration
	err := gormDB.AutoMigrate(models.Product{})
	if err != nil {
		p.FailNowf("error in seed database", err.Error())
	}

	db, err := gormDB.DB()
	if err != nil {
		p.FailNowf("error in seed database", err.Error())
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
		p.FailNowf("error in seed database", err.Error())
	}

	result, err := gormPostgres.Paginate[*models.Product](p.ctx, utils.NewListQuery(10, 1), gormDB)
	items = result.Items
}
