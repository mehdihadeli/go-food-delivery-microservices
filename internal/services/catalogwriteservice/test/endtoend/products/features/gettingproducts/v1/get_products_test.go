//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"

	"github.com/gavv/httpexpect/v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var integrationFixture *integration.IntegrationTestSharedFixture

func TestGetAllProductEndpoint(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "GetAllProducts Endpoint EndToEnd Tests")
}

var _ = Describe("Get All Products Feature", func() {
	var ctx context.Context

	_ = BeforeEach(func() {
		ctx = context.Background()

		By("Seeding the required data")
		integrationFixture.SetupTest()
	})

	_ = AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	// "Scenario" step for testing the get all products API
	Describe("Get all products returns ok status", func() {
		// "When" step
		When("A request is made to get all products", func() {
			// "Then" step
			It("Should return an OK status", func() {
				// Create an HTTPExpect instance and make the request
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				expect.GET("products").
					WithContext(ctx).
					Expect().
					Status(http.StatusOK)
			})
		})
	})
})
