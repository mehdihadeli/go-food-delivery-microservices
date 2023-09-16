//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	uuid "github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var integrationFixture *integration.IntegrationTestSharedFixture

func TestUpdateProductEndpoint(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "UpdateProduct Endpoint EndToEnd Tests")
}

var _ = Describe("UpdateProductE2ETest Suite", func() {
	var (
		ctx     context.Context
		id      uuid.UUID
		request *dtos.UpdateProductRequestDto
	)

	_ = BeforeEach(func() {
		ctx = context.Background()

		By("Seeding the required data")
		integrationFixture.InitializeTest()
		id = integrationFixture.Items[0].ProductId
	})

	_ = AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.DisposeTest()
	})

	// "Scenario" step for testing the update product API with valid input
	Describe("Update product with valid input returns NoContent status", func() {
		BeforeEach(func() {
			request = &dtos.UpdateProductRequestDto{
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       gofakeit.Price(100, 1000),
				Name:        gofakeit.Name(),
			}
		})

		// "When" step
		When("A valid request is made to update a product", func() {
			// "Then" step
			It("Should return a NoContent status", func() {
				// Create an HTTPExpect instance and make the request
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				expect.PUT("products/{id}").
					WithPath("id", id.String()).
					WithJSON(request).
					WithContext(ctx).
					Expect().
					Status(http.StatusNoContent)
			})
		})
	})

	// "Scenario" step for testing the update product API with invalid input
	Describe("Update product returns BadRequest with invalid input", func() {
		BeforeEach(func() {
			// Get a valid product ID from your test data
			id = uuid.NewV4()
			request = &dtos.UpdateProductRequestDto{
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       0,
				Name:        gofakeit.Name(),
			}
		})
		// "When" step
		When("An invalid request is made to update a product", func() {
			// "Then" step
			It("Should return a BadRequest status", func() {
				// Create an HTTPExpect instance and make the request
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				expect.PUT("products/{id}").
					WithPath("id", id.String()).
					WithJSON(request).
					WithContext(context.Background()).
					Expect().
					Status(http.StatusBadRequest)
			})
		})
	})
})
