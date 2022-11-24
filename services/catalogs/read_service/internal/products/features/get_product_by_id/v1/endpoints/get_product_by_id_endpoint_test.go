package endpoints

import (
	"context"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/shared/test_fixture/e2e"
)

// we could also run the server on docker and then send rest call to the api
func Test_Product_By_Id_E2E(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := e2e.NewE2ETestFixture()

	e := NewGetProductByIdEndpoint(delivery.NewProductEndpointBase(fixture.InfrastructureConfigurations, fixture.V1.ProductsGroup, fixture.Bus, fixture.CatalogsMetrics))
	e.MapRoute()

	fixture.Run()
	defer fixture.Cleanup()

	// create httpexpect instance
	expect := httpexpect.New(t, fixture.HttpServer.URL)

	expect.GET("/api/v1/products/{id}").
		WithPath("id", "6b60d642-97ff-4210-baeb-05014d346a48").
		WithContext(context.Background()).
		Expect().
		Status(http.StatusOK)
}
