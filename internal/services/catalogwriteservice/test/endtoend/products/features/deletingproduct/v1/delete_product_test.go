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

func TestDeleteProductEndpoint(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "DeleteProduct Endpoint EndToEnd Tests")
}

var _ = Describe("Delete Product Feature", func() {
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

	// "Scenario" step for testing the delete product API with valid input
	Describe("Delete product with valid input returns NoContent status", func() {
		// "When" step
		When("A valid request is made to delete a product", func() {
			// "Then" step
			It("Should return a NoContent status", func() {
				// Create an HTTPExpect instance and make the request
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				expect.DELETE("products/{id}").
					WithContext(ctx).
					WithPath("id", id.String()).
					Expect().
					Status(http.StatusNoContent)
			})
		})
	})

	// "Scenario" step for testing the delete product API with invalid ID
	Describe("Delete product with with invalid ID returns NotFound status", func() {
		BeforeEach(func() {
			// Generate an invalid UUID
			id = uuid.NewV4()
		})

		// "When" step
		When("An invalid request is made with an invalid ID", func() {
			// "Then" step
			It("Should return a NotFound status", func() {
				// Create an HTTPExpect instance and make the request
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				expect.DELETE("products/{id}").
					WithContext(ctx).
					WithPath("id", id.String()).
					Expect().
					Status(http.StatusNotFound)
			})
		})
	})
})
