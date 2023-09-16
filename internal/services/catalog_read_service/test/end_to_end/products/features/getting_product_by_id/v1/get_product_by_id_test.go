//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/test_fixture/integration"

	"github.com/gavv/httpexpect/v2"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetProductById(t *testing.T) {
	e2eFixture := integration.NewIntegrationTestSharedFixture(t)

	Convey("Get Product By Id Feature", t, func() {
		e2eFixture.InitializeTest()

		ctx := context.Background()
		id := e2eFixture.Items[0].Id

		// "Scenario" step for testing the get product by ID API with a valid ID
		Convey("Get product by ID with a valid ID returns ok status", func() {
			Convey("When A valid request is made with a valid ID", func() {
				expect := httpexpect.New(t, e2eFixture.BaseAddress)

				Convey("Then the response status should be OK", func() {
					expect.GET("products/{id}").
						WithPath("id", id).
						WithContext(ctx).
						Expect().
						Status(http.StatusOK)
				})
			})
		})

		e2eFixture.DisposeTest()
	})
}
