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

func Test_Orders_E2E(t *testing.T) {
	test.SkipCI(t)
	fixture := e2e.NewE2ETestFixture()

	e := NewGetOrdersEndpoint(delivery.NewOrderEndpointBase(fixture.InfrastructureConfiguration, fixture.V1.OrdersGroup))
	e.MapRoute()

	defer fixture.Cleanup()

	s := httptest.NewServer(fixture.Echo)
	defer s.Close()

	// create httpexpect instance
	expect := httpexpect.New(t, s.URL)

	expect.GET("/api/v1/orders").
		WithContext(context.Background()).
		Expect().
		Status(http.StatusOK)
}
