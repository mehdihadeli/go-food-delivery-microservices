//go:build integration
// +build integration

package v1

import (
	"context"
	"fmt"
	"testing"
	"time"

	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/hypothesis"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/test/messaging"
	v1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/deletingproduct/v1"
	integrationEvents "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/deletingproduct/v1/events/integrationevents"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"

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
		command       *v1.DeleteProduct
		result        *mediatr.Unit
		id            uuid.UUID
		notExistsId   uuid.UUID
		shouldPublish hypothesis.Hypothesis[*integrationEvents.ProductDeletedV1]
	)

	_ = BeforeEach(func() {
		By("Seeding the required data")
		// call base SetupTest hook before running child hook
		integrationFixture.SetupTest()

		// child hook codes should be here
		id = integrationFixture.Items[0].Id
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

	// "Scenario" step for testing deleting an existing product
	Describe("Deleting an existing product from the database", func() {
		Context("Given product already exists in the system", func() {
			BeforeEach(func() {
				shouldPublish = messaging.ShouldProduced[*integrationEvents.ProductDeletedV1](
					ctx,
					integrationFixture.Bus,
					nil,
				)
				command, err = v1.NewDeleteProduct(id)
				Expect(err).ShouldNot(HaveOccurred())
			})

			When(
				"the DeleteProduct command is executed for existing product",
				func() {
					BeforeEach(func() {
						result, err = mediatr.Send[*v1.DeleteProduct, *mediatr.Unit](
							ctx,
							command,
						)
					})

					It("Should not return an error", func() {
						Expect(err).NotTo(HaveOccurred())
					})

					It("Should delete the product from the database", func() {
						deletedProduct, err := integrationFixture.CatalogsDBContext.FindProductByID(
							ctx,
							id,
						)
						Expect(err).To(HaveOccurred())
						Expect(
							err,
						).To(MatchError(ContainSubstring(fmt.Sprintf("product with id `%s` not found in the database", command.ProductID.String()))))
						Expect(deletedProduct).To(BeNil())
					})
				},
			)
		})
	})

	// "Scenario" step for testing deleting a non-existing product
	Describe("Deleting a non-existing product from the database", func() {
		Context("Given product does not exists in the system", func() {
			BeforeEach(func() {
				notExistsId = uuid.NewV4()
				command, err = v1.NewDeleteProduct(notExistsId)
				Expect(err).ShouldNot(HaveOccurred())
			})

			When(
				"the DeleteProduct command is executed for non-existing product",
				func() {
					BeforeEach(func() {
						result, err = mediatr.Send[*v1.DeleteProduct, *mediatr.Unit](
							ctx,
							command,
						)
					})

					It("Should return an error", func() {
						Expect(err).To(HaveOccurred())
					})

					It("Should return a NotFound error", func() {
						Expect(
							err,
						).To(MatchError(ContainSubstring(fmt.Sprintf("product with id `%s` not found in the database", command.ProductID.String()))))
					})

					It("Should return a custom NotFound error", func() {
						Expect(customErrors.IsNotFoundError(err)).To(BeTrue())
					})

					It("Should not return a result", func() {
						Expect(result).To(BeNil())
					})
				},
			)
		})
	})

	Describe(
		"Publishing ProductDeleted event when product deleted successfully",
		func() {
			Context("Given product already exists in the system", func() {
				BeforeEach(func() {
					shouldPublish = messaging.ShouldProduced[*integrationEvents.ProductDeletedV1](
						ctx,
						integrationFixture.Bus,
						nil,
					)
					command, err = v1.NewDeleteProduct(id)
					Expect(err).ShouldNot(HaveOccurred())
				})

				When(
					"the DeleteProduct command is executed for existing product",
					func() {
						BeforeEach(func() {
							result, err = mediatr.Send[*v1.DeleteProduct, *mediatr.Unit](
								ctx,
								command,
							)
						})

						It(
							"Should publish ProductDeleted event to the broker",
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
