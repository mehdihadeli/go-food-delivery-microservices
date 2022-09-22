package v1

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/e2e"
	"net/http"
	"testing"
)

// we could also run the server on docker and then send rest call to the api
func Test_Create_Product_E2E(t *testing.T) {
	test.SkipCI(t)
	fixture := e2e.NewE2ETestFixture()

	e := NewCreteProductEndpoint(delivery.NewProductEndpointBase(fixture.InfrastructureConfiguration, fixture.V1.ProductsGroup))
	e.MapRoute()

	fixture.Run()
	defer fixture.Cleanup()

	request := dtos.CreateProductRequestDto{
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(100, 1000),
		Name:        gofakeit.Name(),
	}

	// create httpexpect instance
	expect := httpexpect.New(t, fixture.HttpServer.URL)

	expect.POST("/api/v1/products").
		WithContext(context.Background()).
		WithJSON(request).
		Expect().
		Status(http.StatusCreated)
}
