//go:build integration
// +build integration

package v1

import (
	"context"
	"net/http"
	"testing"
	"time"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_product_by_id/v1/dtos"
	getProductByIdQuery "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_product_by_id/v1/queries"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"

	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var integrationFixture *integration.IntegrationTestSharedFixture

func TestGetProductById(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "Get Product By ID Integration Tests")
}

var _ = Describe("Get Product by ID Feature", func() {
	var (
		ctx    context.Context
		id     uuid.UUID
		query  *getProductByIdQuery.GetProductById
		result *dtos.GetProductByIdResponseDto
		err    error
	)

	_ = BeforeEach(func() {
		By("Seeding the required data")
		integrationFixture.InitializeTest()

		id = integrationFixture.Items[0].ProductId
	})

	_ = AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.DisposeTest()
	})

	_ = BeforeSuite(func() {
		ctx = context.Background()

		// in test mode we set rabbitmq `AutoStart=false` in configuration in rabbitmqOptions, so we should run rabbitmq bus manually
		err = integrationFixture.Bus.Start(context.Background())
		Expect(err).ShouldNot(HaveOccurred())

		// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
		time.Sleep(1 * time.Second)
	})

	_ = AfterSuite(func() {
		integrationFixture.Log.Info("TearDownSuite started")
		err := integrationFixture.Bus.Stop()
		Expect(err).ShouldNot(HaveOccurred())
		time.Sleep(1 * time.Second)
	})

	// "Scenario" step for testing returning an existing product with correct properties
	Describe("Returning an existing product with valid Id from the database with correct properties", func() {
		Context("Given products exists in the database", func() {
			BeforeEach(func() {
				query, err = getProductByIdQuery.NewGetProductById(id)
			})

			// "When" step
			When("the GteProductById query is executed for existing product", func() {
				BeforeEach(func() {
					result, err = mediatr.Send[*getProductByIdQuery.GetProductById, *dtos.GetProductByIdResponseDto](
						ctx,
						query,
					)
				})

				// "Then" step
				It("Should not return an error", func() {
					Expect(err).To(BeNil())
				})

				It("Should return a non-nil result", func() {
					Expect(result).NotTo(BeNil())
				})

				It("Should return a product with the correct ID", func() {
					Expect(result.Product).NotTo(BeNil())
					Expect(result.Product.ProductId).To(Equal(id))
				})
			})
		})
	})

	// "Scenario" step for testing returning a NotFound error when record does not exist
	Describe("Returning a NotFound error when product with specific id does not exist", func() {
		Context("Given products does not exists in the database", func() {
			BeforeEach(func() {
				// Generate a random UUID that does not exist in the database
				id = uuid.NewV4()
				query, err = getProductByIdQuery.NewGetProductById(id)
				Expect(err).To(BeNil())
			})

			// "When" step
			When("the GteProductById query is executed for non-existing product", func() {
				BeforeEach(func() {
					result, err = mediatr.Send[*getProductByIdQuery.GetProductById, *dtos.GetProductByIdResponseDto](
						ctx,
						query,
					)
				})

				// "Then" step
				It("Should return an error", func() {
					Expect(err).To(HaveOccurred())
				})

				It("Should return a NotFound error", func() {
					Expect(err).To(MatchError(ContainSubstring("error in getting product with id")))
				})

				It("Should return a custom NotFound error", func() {
					Expect(customErrors.IsNotFoundError(err)).To(BeTrue())
					Expect(customErrors.IsApplicationError(err, http.StatusNotFound)).To(BeTrue())
				})

				It("Should not return a result", func() {
					Expect(result).To(BeNil())
				})
			})
		})
	})
})
