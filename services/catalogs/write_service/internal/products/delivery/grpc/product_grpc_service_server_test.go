package grpc

import (
	"context"
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
	*ProductGrpcServiceServer
}

func TestRunner(t *testing.T) {
	test.SkipCI(t)
	fixture := e2e.NewE2ETestFixture()

	//https://pkg.go.dev/testing@master#hdr-Subtests_and_Sub_benchmarks
	t.Run("GRPC", func(t *testing.T) {
		// Before running the tests
		productGrpcService := NewProductGrpcService(fixture.InfrastructureConfiguration)
		productService.RegisterProductsServiceServer(fixture.GrpcServer.GetCurrentGrpcServer(), productGrpcService)

		go func() {
			if err := fixture.GrpcServer.RunGrpcServer(nil); err != nil {
				fixture.Log.Errorf("(s.RunGrpcServer) err: %v", err)
			}
		}()

		productGrpcServiceTests := ProductGrpcServiceTests{
			T:                        t,
			E2ETestFixture:           fixture,
			ProductGrpcServiceServer: productGrpcService,
		}

		// Run Tests
		productGrpcServiceTests.Test_GetProduct_By_Id()
		productGrpcServiceTests.Test_Create_Product()

		// After running the tests
		fixture.GrpcServer.GracefulShutdown()
		fixture.Cleanup()
	})
}

func (p *ProductGrpcServiceTests) Test_Create_Product() {
	request := &productService.CreateProductReq{
		Price:       gofakeit.Price(100, 1000),
		Name:        gofakeit.Name(),
		Description: gofakeit.AdjectiveDescriptive(),
	}

	res, err := p.CreateProduct(context.Background(), request)
	assert.NoError(p.T, err)
	assert.NotZero(p.T, res.ProductId)
}

func (p *ProductGrpcServiceTests) Test_GetProduct_By_Id() {
	res, err := p.GetProductById(context.Background(), &productService.GetProductByIdReq{ProductId: "1b088075-53f0-4376-a491-ca6fe3a7f8fa"})
	assert.NoError(p.T, err)
	assert.NotNil(p.T, res.Product)
	assert.Equal(p.T, res.Product.ProductId, "1b088075-53f0-4376-a491-ca6fe3a7f8fa")
}
