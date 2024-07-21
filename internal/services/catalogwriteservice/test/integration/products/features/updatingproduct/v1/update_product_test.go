//go:build integration
// +build integration

package v1

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/hypothesis"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/messaging"
	datamodel "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/data/datamodels"
	v1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/updatingproduct/v1"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/updatingproduct/v1/events/integrationevents"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"

	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var integrationFixture *integration.IntegrationTestSharedFixture

func TestUpdateProduct(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "Updated Products Integration Tests")
}

var _ = Describe("Update Product Feature", func() {
	// Define variables to hold command and result data
	var (
		ctx             context.Context
		existingProduct *datamodel.ProductDataModel
		command         *v1.UpdateProduct
		result          *mediatr.Unit
		err             error
		id              uuid.UUID
		shouldPublish   hypothesis.Hypothesis[*integrationevents.ProductUpdatedV1]
	)

	_ = BeforeEach(func() {
		By("Seeding the required data")
		integrationFixture.SetupTest()

		existingProduct = integrationFixture.Items[0]
	})

	_ = AfterEach(func() {
		By("Cleanup test data")
		integrationFixture.TearDownTest()
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

	// "Scenario" step for testing updating an existing product
	Describe("Updating an existing product in the database", func() {
		Context("Given product exists in the database", func() {
			BeforeEach(func() {
				command, err = v1.NewUpdateProduct(
					existingProduct.Id,
					"Updated Product ShortTypeName",
					existingProduct.Description,
					existingProduct.Price,
				)
				Expect(err).NotTo(HaveOccurred())
			})

			// "When" step
			When("the UpdateProduct command is executed", func() {
				BeforeEach(func() {
					result, err = mediatr.Send[*v1.UpdateProduct, *mediatr.Unit](
						ctx,
						command,
					)
				})

				// "Then" step
				It("Should not return an error", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(result).NotTo(BeNil())
				})

				It("Should return a non-nil result", func() {
					Expect(result).NotTo(BeNil())
				})

				It(
					"Should update the existing product in the database",
					func() {
						updatedProduct, err := integrationFixture.CatalogsDBContext.FindProductByID(
							ctx,
							existingProduct.Id,
						)
						Expect(err).To(BeNil())
						Expect(updatedProduct).NotTo(BeNil())
						Expect(
							updatedProduct.Id,
						).To(Equal(existingProduct.Id))
						Expect(
							updatedProduct.Price,
						).To(Equal(existingProduct.Price))
						Expect(
							updatedProduct.Name,
						).NotTo(Equal(existingProduct.Name))
					},
				)
			})
		})
	})

	// "Scenario" step for testing updating a non-existing product
	Describe("Updating a non-existing product in the database", func() {
		Context("Given product not exists in the database", func() {
			BeforeEach(func() {
				// Generate a random ID that does not exist in the database
				id = uuid.NewV4()
				command, err = v1.NewUpdateProduct(
					id,
					"Updated Product ShortTypeName",
					"Updated Product Description",
					100,
				)
				Expect(err).NotTo(HaveOccurred())
			})

			// "When" step
			When(
				"the UpdateProduct command executed for non-existing product",
				func() {
					BeforeEach(func() {
						result, err = mediatr.Send[*v1.UpdateProduct, *mediatr.Unit](
							ctx,
							command,
						)
					})

					// "Then" step
					It("Should return an error", func() {
						Expect(err).To(HaveOccurred())
					})
					It("Should not return a result", func() {
						Expect(result).To(BeNil())
					})

					It("Should return a NotFound error", func() {
						Expect(
							err,
						).To(MatchError(ContainSubstring(fmt.Sprintf("product with id `%s` not found", id.String()))))
					})

					It("Should return a custom NotFound error", func() {
						Expect(customErrors.IsNotFoundError(err)).To(BeTrue())
						Expect(
							customErrors.IsApplicationError(
								err,
								http.StatusNotFound,
							),
						).To(BeTrue())
					})
				},
			)
		})
	})

	// "Scenario" step for testing updating an existing product
	Describe(
		"Publishing ProductUpdated when product updated  successfully",
		func() {
			Context("Given product exists in the database", func() {
				BeforeEach(func() {
					command, err = v1.NewUpdateProduct(
						existingProduct.Id,
						"Updated Product ShortTypeName",
						existingProduct.Description,
						existingProduct.Price,
					)
					Expect(err).NotTo(HaveOccurred())

					shouldPublish = messaging.ShouldProduced[*integrationevents.ProductUpdatedV1](
						ctx,
						integrationFixture.Bus,
						nil,
					)
				})

				// "When" step
				When(
					"the UpdateProduct command is executed for existing product",
					func() {
						BeforeEach(func() {
							result, err = mediatr.Send[*v1.UpdateProduct, *mediatr.Unit](
								ctx,
								command,
							)
						})

						It(
							"Should publish ProductUpdated event to the broker",
							func() {
								// ensuring message published to the rabbitmq broker
								shouldPublish.Validate(
									ctx,
									"there is no published message",
									time.Second*30,
								)
							},
						)
					},
				)
			})
		},
	)
})
