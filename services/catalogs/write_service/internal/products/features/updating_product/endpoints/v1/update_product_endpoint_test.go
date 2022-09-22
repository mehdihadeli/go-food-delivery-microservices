package v1

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/e2e"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"testing"
)

// we could also run the server on docker and then send rest call to the api
func Test_Update_Product_E2E(t *testing.T) {
	test.SkipCI(t)
	fixture := e2e.NewE2ETestFixture()

	e := NewUpdateProductEndpoint(delivery.NewProductEndpointBase(fixture.InfrastructureConfiguration, fixture.V1.ProductsGroup))
	e.MapRoute()

	fixture.Run()
	defer fixture.Cleanup()

	id, err := uuid.FromString("49a8e487-945b-4050-9a4c-a9242247cb48")
	if err != nil {
		return
	}
	request := updating_product.UpdateProductRequestDto{
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(100, 1000),
		Name:        gofakeit.Name(),
	}

	// create httpexpect instance
	expect := httpexpect.New(t, fixture.HttpServer.URL)

	expect.PUT("/api/v1/products/{id}").
		WithContext(context.Background()).
		WithJSON(request).
		WithPath("id", id.String()).
		Expect().
		Status(http.StatusNoContent)
}
