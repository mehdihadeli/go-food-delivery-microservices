//go:build integration
// +build integration

package v1

import (
	"context"
	"net/http"
	"testing"
	"time"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/hypothesis"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/deleting_product/v1/commands"
	integrationEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/deleting_product/v1/events/integration_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"

	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var integrationFixture *integration.IntegrationTestSharedFixture

func TestDeleteProduct(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "Delete Product Integration Tests")
}

// https://specflow.org/learn/gherkin/#learn-gherkin
// scenario
var _ = Describe("Delete Product Feature", func() {
	var (
		ctx           context.Context
		err           error
		command       *commands.DeleteProduct
		result        *mediatr.Unit
		id            uuid.UUID
		notExistsId   uuid.UUID
		shouldPublish hypothesis.Hypothesis[*integrationEvents.ProductDeletedV1]
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

	// "Scenario" step for testing deleting an existing product
	Describe("Deleting an existing product from the database", func() {
		Context("Given product already exists in the system", func() {
			BeforeEach(func() {
				shouldPublish = messaging.ShouldProduced[*integrationEvents.ProductDeletedV1](
					ctx,
					integrationFixture.Bus,
					nil,
				)
				command, err = commands.NewDeleteProduct(id)
				Expect(err).ShouldNot(HaveOccurred())
			})

			When("the DeleteProduct command is executed for existing product", func() {
				BeforeEach(func() {
					result, err = mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](
						ctx,
						command,
					)
				})

				It("Should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("Should delete the product from the database", func() {
					deletedProduct, err := integrationFixture.ProductRepository.GetProductById(ctx, id)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("can't find the product with id")))
					Expect(deletedProduct).To(BeNil())
				})
			})
		})
	})

	// "Scenario" step for testing deleting a non-existing product
	Describe("Deleting a non-existing product from the database", func() {
		Context("Given product does not exists in the system", func() {
			BeforeEach(func() {
				notExistsId = uuid.NewV4()
				command, err = commands.NewDeleteProduct(notExistsId)
				Expect(err).ShouldNot(HaveOccurred())
			})

			When("the DeleteProduct command is executed for non-existing product", func() {
				BeforeEach(func() {
					result, err = mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](
						ctx,
						command,
					)
				})

				It("Should return an error", func() {
					Expect(err).To(HaveOccurred())
				})

				It("Should return a NotFound error", func() {
					Expect(err).To(MatchError(ContainSubstring("product not found")))
				})

				It("Should return a custom NotFound error", func() {
					Expect(customErrors.IsApplicationError(err, http.StatusNotFound)).To(BeTrue())
					Expect(customErrors.IsNotFoundError(err)).To(BeTrue())
				})

				It("Should not return a result", func() {
					Expect(result).To(BeNil())
				})
			})
		})
	})

	Describe("Publishing ProductDeleted event when product deleted successfully", func() {
		Context("Given product already exists in the system", func() {
			BeforeEach(func() {
				shouldPublish = messaging.ShouldProduced[*integrationEvents.ProductDeletedV1](
					ctx,
					integrationFixture.Bus,
					nil,
				)
				command, err = commands.NewDeleteProduct(id)
				Expect(err).ShouldNot(HaveOccurred())
			})

			When("the DeleteProduct command is executed for existing product", func() {
				BeforeEach(func() {
					result, err = mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](
						ctx,
						command,
					)
				})

				It("Should publish ProductDeleted event to the broker", func() {
					// ensuring message published to the rabbitmq broker
					shouldPublish.Validate(ctx, "there is no published message", time.Second*30)
				})
			})
		})
	})
})
