//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogreadservice/internal/shared/testfixture/integration"

	"github.com/gavv/httpexpect/v2"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAllProducts(t *testing.T) {
	e2eFixture := integration.NewIntegrationTestSharedFixture(t)

	Convey("Get All Products Feature", t, func() {
		e2eFixture.SetupTest()
		ctx := context.Background()

		Convey("Get all products returns ok status", func() {
			Convey("When a request is made to get all products", func() {
				expect := httpexpect.New(t, e2eFixture.BaseAddress)

				Convey("Then the response status should be OK", func() {
					expect.GET("products").
						WithContext(ctx).
						Expect().
						Status(http.StatusOK)
				})
			})
		})

		e2eFixture.TearDownTest()
	})
}
