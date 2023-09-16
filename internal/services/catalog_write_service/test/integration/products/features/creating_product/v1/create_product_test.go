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
	createProductCommand "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/dtos"
	integrationEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/events/integration_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var integrationFixture *integration.IntegrationTestSharedFixture

func TestCreateProduct(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "Create Product Integration Tests")
}

// https://specflow.org/learn/gherkin/#learn-gherkin
// scenario
var _ = Describe("Creating Product Feature", func() {
	var (
		ctx            context.Context
		err            error
		command        *createProductCommand.CreateProduct
		result         *dtos.CreateProductResponseDto
		createdProduct *models.Product
		id             uuid.UUID
		shouldPublish  hypothesis.Hypothesis[*integrationEvents.ProductCreatedV1]
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

	// "Scenario" step for testing creating a new product
	Describe("Creating a new product and saving it to the database when product doesn't exists", func() {
		Context("Given new product doesn't exists in the system", func() {
			BeforeEach(func() {
				command, err = createProductCommand.NewCreateProduct(
					gofakeit.Name(),
					gofakeit.AdjectiveDescriptive(),
					gofakeit.Price(150, 6000),
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(command).ToNot(BeNil())
			})

			When("the CreateProduct command is executed for non-existing product", func() {
				BeforeEach(func() {
					result, err = mediatr.Send[*createProductCommand.CreateProduct, *dtos.CreateProductResponseDto](
						ctx,
						command,
					)
				})

				It("Should create the product successfully", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(result).NotTo(BeNil())
				})

				It("Should have a non-empty product ID matching the command ID", func() {
					Expect(result.ProductID).To(Equal(command.ProductID))
				})

				It("Should be able to retrieve the product from the database", func() {
					createdProduct, err = integrationFixture.ProductRepository.GetProductById(
						ctx,
						result.ProductID,
					)
					Expect(err).NotTo(HaveOccurred())

					Expect(result).NotTo(BeNil())
					Expect(command.ProductID).To(Equal(result.ProductID))
					Expect(createdProduct).NotTo(BeNil())
				})
			})
		})
	})

	// "Scenario" step for testing creating a product with duplicate data
	Describe("Creating a new product with duplicate data and already exists product", func() {
		Context("Given product already exists in the system", func() {
			BeforeEach(func() {
				command = &createProductCommand.CreateProduct{
					Name:        gofakeit.Name(),
					Description: gofakeit.AdjectiveDescriptive(),
					Price:       gofakeit.Price(150, 6000),
					ProductID:   id,
				}
			})

			When("the CreateProduct command is executed for existing product", func() {
				BeforeEach(func() {
					result, err = mediatr.Send[*createProductCommand.CreateProduct, *dtos.CreateProductResponseDto](
						ctx,
						command,
					)
				})

				It("Should return an error indicating duplicate record", func() {
					Expect(err).To(HaveOccurred())
					Expect(customErrors.IsApplicationError(err, http.StatusConflict)).To(BeTrue())
				})

				It("Should not return a result", func() {
					Expect(result).To(BeNil())
				})
			})
		})
	})

	// "Scenario" step for testing creating a product with duplicate data
	Describe("Publishing ProductCreated event to the broker when product saved successfully", func() {
		Context("Given new product doesn't exists in the system", func() {
			BeforeEach(func() {
				shouldPublish = messaging.ShouldProduced[*integrationEvents.ProductCreatedV1](
					ctx,
					integrationFixture.Bus,
					nil,
				)
				command, err = createProductCommand.NewCreateProduct(
					gofakeit.Name(),
					gofakeit.AdjectiveDescriptive(),
					gofakeit.Price(150, 6000),
				)
				Expect(err).ToNot(HaveOccurred())
			})

			When("CreateProduct command is executed for non-existing product", func() {
				BeforeEach(func() {
					result, err = mediatr.Send[*createProductCommand.CreateProduct, *dtos.CreateProductResponseDto](
						ctx,
						command,
					)
				})

				It("Should return no error", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("Should return not nil result", func() {
					Expect(result).ToNot(BeNil())
				})

				It("Should publish ProductCreated event to the broker", func() {
					// ensuring message published to the rabbitmq broker
					shouldPublish.Validate(ctx, "there is no published message", time.Second*30)
				})
			})
		})
	})
})
