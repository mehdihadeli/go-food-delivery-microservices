package grpc

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/suite"

	productService "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/mocks/testData"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/e2e"
)

type productGrpcServiceE2eTests struct {
	*e2e.E2ETestFixture
	*e2e.E2ETestSharedFixture
	productsServiceClient productService.ProductsServiceClient
}

func TestProductGrpcServiceE2E(t *testing.T) {
	suite.Run(t, &productGrpcServiceE2eTests{E2ETestSharedFixture: e2e.NewE2ETestSharedFixture(t)})
}

func (c *productGrpcServiceE2eTests) Test_Should_Create_Product_With_Valid_Data_In_DB() {
	request := &productService.CreateProductReq{
		Price:       gofakeit.Price(100, 1000),
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
	}

	res, err := c.productsServiceClient.CreateProduct(c.Ctx, request)
	c.NoError(err)
	c.NotEmpty(res.ProductId)
}

func (c *productGrpcServiceE2eTests) Test_Should_Return_Data_With_Valid_Id() {
	id := testData.Products[0].ProductId.String()

	res, err := c.productsServiceClient.GetProductById(c.Ctx, &productService.GetProductByIdReq{ProductId: id})

	fmt.Println(err)
	fmt.Println(res)
	c.NoError(err)
	c.NotNil(res.Product)
	c.Equal(res.Product.ProductId, id)
}

func (c *productGrpcServiceE2eTests) SetupTest() {
	c.T().Log("SetupTest")
	c.E2ETestFixture = e2e.NewE2ETestFixture(c.E2ETestSharedFixture)

	// Before running the tests
	productGrpcService := NewProductGrpcService(c.InfrastructureConfigurations, c.CatalogsMetrics, c.Bus)
	productService.RegisterProductsServiceServer(c.GrpcServer.GetCurrentGrpcServer(), productGrpcService)

	c.E2ETestFixture.Run()

	c.productsServiceClient = productService.NewProductsServiceClient(c.GrpcClient.GetGrpcConnection())
}

func (c *productGrpcServiceE2eTests) TearDownTest() {
	c.T().Log("TearDownTest")
	// cleanup test containers with their hooks
}

func (c *productGrpcServiceE2eTests) SetupSuite() {
}
