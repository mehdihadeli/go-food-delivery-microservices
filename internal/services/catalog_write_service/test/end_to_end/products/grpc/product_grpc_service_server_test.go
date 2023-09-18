//go:build e2e
// +build e2e

package grpc

import (
	"context"
	"testing"

	productService "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/grpc/genproto"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var integrationFixture *integration.IntegrationTestSharedFixture

func TestProductGrpcServiceEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "ProductGrpcService EndToEnd Tests")
}

var _ = Describe("Product Grpc Service Feature", func() {
	var (
		ctx context.Context
		id  uuid.UUID
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

	// "Scenario" step for testing the creation of a product with valid data in the database
	Describe("Creation of a product with valid data in the database", func() {
		// "When" step
		When("A request is made to create a product with valid data", func() {
			// "Then" step
			It("Should return a non-empty ProductId", func() {
				// Create a gRPC request with valid data
				request := &productService.CreateProductReq{
					Price:       gofakeit.Price(100, 1000),
					Name:        gofakeit.Name(),
					Description: gofakeit.AdjectiveDescriptive(),
				}

				// Make the gRPC request to create the product
				res, err := integrationFixture.ProductServiceClient.CreateProduct(ctx, request)
				Expect(err).To(BeNil())
				Expect(res).NotTo(BeNil())
				Expect(res.ProductId).NotTo(BeEmpty())
			})
		})
	})

	// "Scenario" step for testing the retrieval of data with a valid ID
	Describe("Retrieve product with a valid ID", func() {
		// "When" step
		When("A request is made to retrieve data with a valid ID", func() {
			// "Then" step
			It("Should return data with a matching ProductId", func() {
				// Make the gRPC request to retrieve data by ID
				res, err := integrationFixture.ProductServiceClient.GetProductById(
					ctx,
					&productService.GetProductByIdReq{ProductId: id.String()},
				)

				Expect(err).To(BeNil())
				Expect(res).NotTo(BeNil())
				Expect(res.Product).NotTo(BeNil())
				Expect(res.Product.ProductId).To(Equal(id.String()))
			})
		})
	})
})
