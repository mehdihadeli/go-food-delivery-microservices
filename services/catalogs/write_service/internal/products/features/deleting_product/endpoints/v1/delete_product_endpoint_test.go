package v1

import (
	"context"
	"github.com/gavv/httpexpect/v2"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/test_fixtures/e2e"
	"net/http"
	"testing"
)

// we could also run the server on docker and then send rest call to the api
func Test_Delete_Product_E2E(t *testing.T) {
	test.SkipCI(t)
	fixture := e2e.NewE2ETestFixture()

	e := NewDeleteProductEndpoint(delivery.NewProductEndpointBase(fixture.InfrastructureConfiguration, fixture.V1.ProductsGroup))
	e.MapRoute()

	fixture.Run()
	defer fixture.Cleanup()

	// create httpexpect instance
	expect := httpexpect.New(t, fixture.HttpServer.URL)

	expect.DELETE("/api/v1/products/{id}").
		WithContext(context.Background()).
		WithPath("id", "d1f2f59e-48dd-456d-bbfd-cb51b566b08c").
		Expect().
		Status(http.StatusNoContent)
}
