//go:build.sh e2e
// +build.sh e2e

package endpoints

import (
	"context"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"

	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/products/delivery"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/read_service/internal/shared/test_fixture/e2e"
)

// we could also run the server on docker and then send rest call to the api
func Test_Product_By_Id_E2E(t *testing.T) {
	testUtils.SkipCI(t)
	fixture := e2e.NewE2ETestFixture(e2e.NewE2ETestSharedFixture(t))

	e := NewGetProductByIdEndpoint(
		delivery.NewProductEndpointBase(
			fixture.InfrastructureConfigurations,
			fixture.ProductsGroup,
			fixture.Bus,
			fixture.CatalogsMetrics,
		),
	)
	e.MapRoute()

	fixture.Run()

	// create httpexpect instance
	expect := httpexpect.New(t, fixture.HttpServer.GetEchoInstance().ListenerAddr().String())

	expect.GET("/api/v1/products/{id}").
		WithPath("id", "6b60d642-97ff-4210-baeb-05014d346a48").
		WithContext(context.Background()).
		Expect().
		Status(http.StatusOK)
}
