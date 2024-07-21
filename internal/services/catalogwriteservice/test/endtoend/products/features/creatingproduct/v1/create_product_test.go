//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/dtos"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var integrationFixture *integration.IntegrationTestSharedFixture

func TestCreateProductEndpoint(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "CreateProduct Endpoint EndToEnd Tests")
}

var _ = Describe("CreateProduct Feature", func() {
	var (
		ctx     context.Context
		request *dtos.CreateProductRequestDto
	)

	_ = BeforeEach(func() {
		ctx = context.Background()

		By("Seeding the required data")
		integrationFixture.SetupTest()
	})

	_ = AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
	})

	// "Scenario" step for testing the create product API with valid input
	Describe("Create new product return created status with valid input", func() {
		BeforeEach(func() {
			// Generate a valid request
			request = &dtos.CreateProductRequestDto{
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       gofakeit.Price(100, 1000),
				Name:        gofakeit.Name(),
			}
		})
		// "When" step
		When("A valid request is made to create a product", func() {
			// "Then" step
			It("Should returns a StatusCreated response", func() {
				// Create an HTTPExpect instance and make the request
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				expect.POST("products").
					WithContext(ctx).
					WithJSON(request).
					Expect().
					Status(http.StatusCreated)
			})
		})
	})

	// "Scenario" step for testing the create product API with invalid price input
	Describe("Create product returns a BadRequest status with invalid price input", func() {
		BeforeEach(func() {
			// Generate an invalid request with zero price
			request = &dtos.CreateProductRequestDto{
				Description: gofakeit.AdjectiveDescriptive(),
				Price:       0.0,
				Name:        gofakeit.Name(),
			}
		})
		// "When" step
		When("An invalid request is made with a zero price", func() {
			// "Then" step
			It("Should return a BadRequest status", func() {
				// Create an HTTPExpect instance and make the request
				expect := httpexpect.New(GinkgoT(), integrationFixture.BaseAddress)
				expect.POST("products").
					WithContext(ctx).
					WithJSON(request).
					Expect().
					Status(http.StatusBadRequest)
			})
		})
	})
})
