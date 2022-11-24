package repositories

import (
	"context"
	"testing"
	"time"

	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	defaultLogger "github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/default_logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb"
	repository2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/mongodb/repository"
	mongo2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/containers/testcontainer/mongo"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/models"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var items []*models.Product

const (
	DatabaseName   = "catalogs"
	CollectionName = "products"
)

// Define the custom testify suite
type productMongoRepositoryTestSuite struct {
	suite.Suite
	productRepository contracts.ProductRepository
	ctx               context.Context
}

func TestProductPostgresRepositoryTestSuite(t *testing.T) {
	suite.Run(t, &productMongoRepositoryTestSuite{})
}

func (p *productMongoRepositoryTestSuite) Test_Create_Product_Should_Create_New_Product_In_DB() {
	ctx := p.ctx

	product := &models.Product{
		Id:          uuid.NewV4().String(),
		ProductId:   uuid.NewV4().String(),
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(100, 1000),
		CreatedAt:   time.Now(),
	}

	createdProduct, err := p.productRepository.CreateProduct(ctx, product)
	require.NoError(p.T(), err)

	product, err = p.productRepository.GetProductById(ctx, createdProduct.Id)
	p.NoError(err)

	p.NotNil(p)
	p.Equal(product.Id, createdProduct.Id)
}

func (p *productMongoRepositoryTestSuite) Test_Update_Product_Should_Update_Existing_Product_In_DB() {
	ctx := p.ctx

	id := items[0].Id
	existingProduct, err := p.productRepository.GetProductById(ctx, id)
	p.Require().NoError(err)
	p.Require().NotNil(existingProduct)

	existingProduct.Name = "test_update_product"
	_, err = p.productRepository.UpdateProduct(ctx, existingProduct)
	p.Require().NoError(err)

	updatedProduct, err := p.productRepository.GetProductById(ctx, id)
	p.Equal(updatedProduct.Name, "test_update_product")
}

func (p *productMongoRepositoryTestSuite) Test_Delete_Product_Should_Delete_Existing_Product_In_DB() {
	ctx := p.ctx

	id := items[0].Id

	err := p.productRepository.DeleteProductByID(ctx, id)
	p.Require().NoError(err)

	product, err := p.productRepository.GetProductById(ctx, id)

	p.Error(err)
	p.True(customErrors.IsNotFoundError(err))
	p.Nil(product)
}

func (p *productMongoRepositoryTestSuite) Test_Get_Product() {
	ctx := p.ctx
	id := items[0].Id

	p.Run("Should_Return_NotFound_Error_When_Item_DoesNot_Exists", func() {
		// with subset test a new t will create for subset test
		res, err := p.productRepository.GetProductById(ctx, uuid.NewV4().String())

		p.Error(err)
		p.True(customErrors.IsNotFoundError(err))
		p.Nil(res)
	})

	p.Run("Should_Get_Existing_Product_From_DB", func() {
		res, err := p.productRepository.GetProductById(ctx, id)
		p.Require().NoError(err)

		p.NotNil(res)
		p.Equal(res.Id, id)
	})

	p.Run("Should_Get_All_Existing_Products_From_DB", func() {
		res, err := p.productRepository.GetAllProducts(ctx, utils.NewListQuery(10, 1))
		p.Require().NoError(err)

		p.Equal(2, len(res.Items))
	})
}

func (p *productMongoRepositoryTestSuite) SetupSuite() {
	p.T().Log("SetupSuite")
}

func (p *productMongoRepositoryTestSuite) SetupTest() {
	p.ctx = context.Background()
	p.T().Log("SetupTest")

	rep, err := setupTest(p.ctx, p)
	if err != nil {
		p.FailNowf("error in the setup repository", err.Error())
	}

	p.productRepository = rep
}

func (p *productMongoRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	p.T().Log("BeforeTest")
}

func (p *productMongoRepositoryTestSuite) AfterTest(suiteName, testName string) {
	p.T().Log("AfterTest")
}

func (p *productMongoRepositoryTestSuite) TearDownSuite() {
	p.T().Log("TearDownSuite")
}

func (p *productMongoRepositoryTestSuite) TearDownTest() {
	p.T().Log("TearDownTest")
	// cleanup test containers
	p.ctx.Done()
}

func setupTest(ctx context.Context, p *productMongoRepositoryTestSuite) (contracts.ProductRepository, error) {
	mongoDB, err := mongo2.NewMongoTestContainers().Start(ctx, p.T())
	if err != nil {
		return nil, err
	}

	seedAndMigration(p, mongoDB)

	if err != nil {
		return nil, err
	}

	genericRepository := repository2.NewGenericMongoRepository[*models.Product](mongoDB, DatabaseName, CollectionName)
	productRepository := NewMongoProductRepository(defaultLogger.Logger, genericRepository)

	return productRepository, nil
}

func seedAndMigration(p *productMongoRepositoryTestSuite, db *mongo.Client) {
	// https://github.com/go-testfixtures/testfixtures#templating
	// seed data
	seedProducts := []models.Product{
		{Id: uuid.NewV4().String(), ProductId: uuid.NewV4().String(), Name: gofakeit.Name(), Description: gofakeit.AdjectiveDescriptive(), Price: gofakeit.Price(150, 6000), CreatedAt: time.Now()},
		{Id: uuid.NewV4().String(), ProductId: uuid.NewV4().String(), Name: gofakeit.Name(), Description: gofakeit.AdjectiveDescriptive(), Price: gofakeit.Price(150, 6000), CreatedAt: time.Now()},
	}

	//// https://go.dev/doc/faq#convert_slice_of_interface
	data := make([]interface{}, len(seedProducts))
	for i, v := range seedProducts {
		data[i] = v
	}

	collection := db.Database(DatabaseName).Collection(CollectionName)
	_, err := collection.InsertMany(context.Background(), data, &options.InsertManyOptions{})
	if err != nil {
		p.FailNowf("error in seed database", err.Error())
	}

	result, err := mongodb.Paginate[*models.Product](p.ctx, utils.NewListQuery(10, 1), collection, nil)
	items = result.Items
}
