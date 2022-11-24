package endpoints

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/test_fixture/e2e"
)

// we could also run the server on docker and then send rest call to the api
func Test_Get_All_Products_E2E(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := e2e.NewE2ETestFixture()

	e := NewGetProductsEndpoint(delivery.NewProductEndpointBase(fixture.InfrastructureConfigurations, fixture.V1.ProductsGroup, fixture.Bus, fixture.CatalogsMetrics))
	e.MapRoute()

	fixture.Run()
	defer fixture.Cleanup()

	// create httpexpect instance
	expect := httpexpect.New(t, fixture.HttpServer.URL)

	expect.GET("/api/v1/products").
		WithContext(fixture.Ctx).
		Expect().
		Status(http.StatusOK)
}
