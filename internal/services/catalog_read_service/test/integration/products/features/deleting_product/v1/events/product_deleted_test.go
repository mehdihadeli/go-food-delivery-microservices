//go:build integration
// +build integration

package events

import (
	"context"
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
	externalEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/deleting_products/v1/events/integration_events/external_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/test_fixture/integration"

	uuid "github.com/satori/go.uuid"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProductDeleted(t *testing.T) {
	// Setup and initialization code here.
	integrationTestSharedFixture := integration.NewIntegrationTestSharedFixture(t)
	// in test mode we set rabbitmq `AutoStart=false` in configuration in rabbitmqOptions, so we should run rabbitmq bus manually
	integrationTestSharedFixture.Bus.Start(context.Background())
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)

	Convey("Product Deleted Feature", t, func() {
		ctx := context.Background()
		// will execute with each subtest
		integrationTestSharedFixture.InitializeTest()

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey("Consume ProductDeleted event by consumer", func() {
			event := &externalEvents.ProductDeletedV1{
				Message:   types.NewMessage(uuid.NewV4().String()),
				ProductId: integrationTestSharedFixture.Items[0].ProductId,
			}
			// check for consuming `ProductDeletedV1` message with existing consumer
			hypothesis := messaging.ShouldConsume[*externalEvents.ProductDeletedV1](
				ctx,
				integrationTestSharedFixture.Bus,
				nil,
			)

			Convey("When a ProductDeleted event consumed", func() {
				err := integrationTestSharedFixture.Bus.PublishMessage(
					ctx,
					event,
					nil,
				)
				So(err, ShouldBeNil)

				Convey("Then it should consume the ProductDeleted event", func() {
					hypothesis.Validate(ctx, "there is no consumed message", 30*time.Second)
				})
			})
		})

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey("Delete product in mongo database when a ProductDeleted event consumed", func() {
			event := &externalEvents.ProductDeletedV1{
				Message:   types.NewMessage(uuid.NewV4().String()),
				ProductId: integrationTestSharedFixture.Items[0].ProductId,
			}

			Convey("When a ProductDeleted event consumed", func() {
				err := integrationTestSharedFixture.Bus.PublishMessage(
					ctx,
					event,
					nil,
				)
				So(err, ShouldBeNil)

				Convey("It should delete product in the mongo database", func() {
					ctx := context.Background()

					productDeleted := &externalEvents.ProductDeletedV1{
						Message:   types.NewMessage(uuid.NewV4().String()),
						ProductId: integrationTestSharedFixture.Items[0].ProductId,
					}

					err := integrationTestSharedFixture.Bus.PublishMessage(ctx, productDeleted, nil)
					So(err, ShouldBeNil)

					var p *models.Product

					So(testUtils.WaitUntilConditionMet(func() bool {
						p, err = integrationTestSharedFixture.ProductRepository.GetProductByProductId(
							ctx,
							integrationTestSharedFixture.Items[0].ProductId,
						)
						So(err, ShouldBeNil)

						return p == nil
					}), ShouldBeNil)

					So(p, ShouldBeNil)
				})
			})
		})

		integrationTestSharedFixture.DisposeTest()
	})

	integrationTestSharedFixture.Log.Info("TearDownSuite started")
	integrationTestSharedFixture.Bus.Stop()
	time.Sleep(1 * time.Second)
}
