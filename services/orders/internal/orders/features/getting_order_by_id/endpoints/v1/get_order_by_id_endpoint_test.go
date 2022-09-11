package v1

import (
	"context"
	"github.com/gavv/httpexpect/v2"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/test_fixtures/e2e"
	"net/http"
	"net/http/httptest"
	"testing"
)

// we could also run the server on docker and then send rest call to the api
func Test_Order_By_Id_E2E(t *testing.T) {
	test.SkipCI(t)
	fixture := e2e.NewE2ETestFixture()

	e := NewGetOrderByIdEndpoint(delivery.NewOrderEndpointBase(fixture.InfrastructureConfiguration, fixture.V1.OrdersGroup))
	e.MapRoute()

	defer fixture.Cleanup()

	s := httptest.NewServer(fixture.Echo)
	defer s.Close()

	// create httpexpect instance
	expect := httpexpect.New(t, s.URL)

	expect.GET("/api/v1/orders/{id}").
		WithPath("id", "97e2d953-ed25-4afb-8578-782cc5d365ba").
		WithContext(context.Background()).
		Expect().
		Status(http.StatusOK)
}
