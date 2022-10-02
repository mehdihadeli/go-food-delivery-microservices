package grpc

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	productService "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/proto/service_clients"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/e2e"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ProductGrpcServiceTests struct {
	*testing.T
	*e2e.E2ETestFixture
	productService.ProductsServiceClient
}

func TestRunner(t *testing.T) {
	test.SkipCI(t)
	fixture := e2e.NewE2ETestFixture()

	//https://pkg.go.dev/testing@master#hdr-Subtests_and_Sub_benchmarks
	t.Run("GRPC", func(t *testing.T) {
		// Before running the tests
		productGrpcService := NewProductGrpcService(fixture.InfrastructureConfiguration)
		productService.RegisterProductsServiceServer(fixture.GrpcServer.GetCurrentGrpcServer(), productGrpcService)

		ctx := fixture.Ctx
		fixture.Run()

		productGrpcClient := productService.NewProductsServiceClient(fixture.GrpcClient.GetGrpcConnection())

		productGrpcServiceTests := ProductGrpcServiceTests{
			T:                     t,
			E2ETestFixture:        fixture,
			ProductsServiceClient: productGrpcClient,
		}

		// Run Tests
		productGrpcServiceTests.Test_GetProduct_By_Id(ctx)
		productGrpcServiceTests.Test_Create_Product(ctx)

		// After running the tests
		fixture.Cleanup()
	})
}

func (p *ProductGrpcServiceTests) Test_Create_Product(ctx context.Context) {
	request := &productService.CreateProductReq{
		Price:       gofakeit.Price(100, 1000),
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
	}

	res, err := p.CreateProduct(ctx, request)
	assert.NoError(p.T, err)
	assert.NotZero(p.T, res.ProductId)
}

func (p *ProductGrpcServiceTests) Test_GetProduct_By_Id(ctx context.Context) {
	res, err := p.GetProductById(ctx, &productService.GetProductByIdReq{ProductId: "1b088075-53f0-4376-a491-ca6fe3a7f8fa"})
	fmt.Println(err)
	fmt.Println(res)
	assert.NoError(p.T, err)
	assert.NotNil(p.T, res.Product)
	assert.Equal(p.T, res.Product.ProductId, "1b088075-53f0-4376-a491-ca6fe3a7f8fa")
}
