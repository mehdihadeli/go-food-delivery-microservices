//go:build integration
// +build integration

package v1

import (
	"context"
	"testing"
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"
	gettingproductsv1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/gettingproducts/v1"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/gettingproducts/v1/dtos"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"

	"github.com/mehdihadeli/go-mediatr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var integrationFixture *integration.IntegrationTestSharedFixture

func TestGetProducts(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "Get Products Integration Tests")
}

var _ = Describe("Get All Products Feature", func() {
	// Define variables to hold query and result data
	var (
		ctx         context.Context
		query       *gettingproductsv1.GetProducts
		queryResult *dtos.GetProductsResponseDto
		err         error
	)

	_ = BeforeEach(func() {
		By("Seeding the required data")
		// call base SetupTest hook before running child hook
		integrationFixture.SetupTest()

		// child hook codes should be here
	})

	_ = AfterEach(func() {
		By("Cleanup test data")
		// call base TearDownTest hook before running child hook
		integrationFixture.TearDownTest()

		// child hook codes should be here
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

	// "Scenario" step for testing getting a list of existing products
	Describe("Getting a list of existing products from the database", func() {
		Context("Given existing products in the database", func() {
			BeforeEach(func() {
				// Create a query to retrieve a list of products
				query, err = gettingproductsv1.NewGetProducts(
					utils.NewListQuery(10, 1),
				)
				Expect(err).To(BeNil())
			})

			// "When" step
			When(
				"the GteProducts query is executed for existing products",
				func() {
					BeforeEach(func() {
						queryResult, err = mediatr.Send[*gettingproductsv1.GetProducts, *dtos.GetProductsResponseDto](
							ctx,
							query,
						)
					})

					// "Then" step
					It("Should not return an error", func() {
						Expect(err).To(BeNil())
					})

					It("Should return a non-nil result", func() {
						Expect(queryResult).NotTo(BeNil())
					})

					It("Should return a list of products with items", func() {
						Expect(queryResult.Products).NotTo(BeNil())
						Expect(queryResult.Products.Items).NotTo(BeEmpty())
					})

					It("Should return the expected number of products", func() {
						// Replace 'len(c.Products)' with the expected number of products
						Expect(
							len(queryResult.Products.Items),
						).To(Equal(len(integrationFixture.Items)))
					})
				},
			)
		})
	})
})
