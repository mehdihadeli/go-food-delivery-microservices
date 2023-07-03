package endpoints

import (
	"testing"
)

func Test_Orders_E2E(t *testing.T) {
	//testUtils.SkipCI(t)
	//fixture := e2e.NewE2ETestFixture()
	//
	//e := NewGetOrdersEndpoint(delivery.NewOrderEndpointBase(fixture.InfrastructureConfigurations, fixture.V1.OrdersGroup, fixture.Bus, fixture.OrdersMetrics))
	//e.MapRoute()
	//
	//fixture.Run()
	//defer fixture.Cleanup()
	//
	//s := httptest.NewServer(fixture.Echo)
	//defer s.Close()
	//
	//// create httpexpect instance
	//expect := httpexpect.New(t, s.URL)
	//
	//expect.GET("/api/v1/orders").
	//	WithContext(context.Background()).
	//	Expect().
	//	Status(http.StatusOK)
}
