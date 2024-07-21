//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"

	"github.com/gavv/httpexpect/v2"
	uuid "github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var integrationFixture *integration.IntegrationTestSharedFixture

func TestGetProductByIdEndpoint(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "GetProductById Endpoint EndToEnd Tests")
}

var _ = Describe("Get Product By Id Feature", func() {
	var (
		ctx context.Context
		id  uuid.UUID
	)

	_ = BeforeEach(func() {
		ctx = context.Background()

		By("Seeding the required data")
		integrationFixture.SetupTest()

		id = integrationFixture.Items[0].Id
	})

	_ = AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	// "Scenario" step for testing the get product by ID API with a valid ID
	Describe("Get product by ID with a valid ID returns ok status", func() {
		// "When" step
		When("A valid request is made with a valid ID", func() {
			// "Then" step
			It("Should return an OK status", func() {
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				expect.GET("products/{id}").
					WithPath("id", id).
					WithContext(ctx).
					Expect().
					Status(http.StatusOK)
			})
		})
	})

	// "Scenario" step for testing the get product by ID API with a valid ID
	Describe("Get product by ID with a invalid ID returns NotFound status", func() {
		BeforeEach(func() {
			// Generate an invalid UUID
			id = uuid.NewV4()
		})
		When("An invalid request is made with an invalid ID", func() {
			// "Then" step
			It("Should return a NotFound status", func() {
				// Create an HTTPExpect instance and make the request
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				expect.GET("products/{id}").
					WithPath("id", id.String()).
					WithContext(ctx).
					Expect().
					Status(http.StatusNotFound)
			})
		})
	})
})
