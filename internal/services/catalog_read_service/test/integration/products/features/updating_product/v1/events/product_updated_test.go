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
	externalEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/updating_products/v1/events/integration_events/external_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/test_fixture/integration"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProductUpdatedConsumer(t *testing.T) {
	// Setup and initialization code here.
	integrationTestSharedFixture := integration.NewIntegrationTestSharedFixture(t)
	// in test mode we set rabbitmq `AutoStart=false` in configuration in rabbitmqOptions, so we should run rabbitmq bus manually
	integrationTestSharedFixture.Bus.Start(context.Background())
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)

	Convey("Product Created Feature", t, func() {
		ctx := context.Background()
		integrationTestSharedFixture.InitializeTest()

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey("Consume ProductUpdated event by consumer", func() {
			hypothesis := messaging.ShouldConsume[*externalEvents.ProductUpdatedV1](
				ctx,
				integrationTestSharedFixture.Bus,
				nil,
			)

			fakeUpdateProduct := &externalEvents.ProductUpdatedV1{
				Message:     types.NewMessage(uuid.NewV4().String()),
				ProductId:   integrationTestSharedFixture.Items[0].ProductId,
				Name:        gofakeit.Name(),
				Price:       gofakeit.Price(100, 1000),
				Description: gofakeit.EmojiDescription(),
				UpdatedAt:   time.Now(),
			}

			Convey("When a ProductUpdated event consumed", func() {
				err := integrationTestSharedFixture.Bus.PublishMessage(ctx, fakeUpdateProduct, nil)
				So(err, ShouldBeNil)

				Convey("Then it should consume the ProductUpdated event", func() {
					hypothesis.Validate(ctx, "there is no consumed message", 30*time.Second)
				})
			})
		})

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey("Update product in mongo database when a ProductDeleted event consumed", func() {
			fakeUpdateProduct := &externalEvents.ProductUpdatedV1{
				Message:     types.NewMessage(uuid.NewV4().String()),
				ProductId:   integrationTestSharedFixture.Items[0].ProductId,
				Name:        gofakeit.Name(),
				Price:       gofakeit.Price(100, 1000),
				Description: gofakeit.EmojiDescription(),
				UpdatedAt:   time.Now(),
			}

			Convey("When a ProductUpdated event consumed", func() {
				err := integrationTestSharedFixture.Bus.PublishMessage(ctx, fakeUpdateProduct, nil)
				So(err, ShouldBeNil)

				Convey("It should update product in the mongo database", func() {
					ctx := context.Background()
					productUpdated := &externalEvents.ProductUpdatedV1{
						Message:     types.NewMessage(uuid.NewV4().String()),
						ProductId:   integrationTestSharedFixture.Items[0].ProductId,
						Name:        gofakeit.Name(),
						Description: gofakeit.AdjectiveDescriptive(),
						Price:       gofakeit.Price(150, 6000),
						UpdatedAt:   time.Now(),
					}

					err := integrationTestSharedFixture.Bus.PublishMessage(ctx, productUpdated, nil)
					So(err, ShouldBeNil)

					var product *models.Product

					err = testUtils.WaitUntilConditionMet(func() bool {
						product, err = integrationTestSharedFixture.ProductRepository.GetProductByProductId(
							ctx,
							integrationTestSharedFixture.Items[0].ProductId,
						)

						return product != nil && product.Name == productUpdated.Name
					})

					So(err, ShouldBeNil)
					So(product, ShouldNotBeNil)
					So(productUpdated.ProductId, ShouldEqual, product.ProductId)
				})
			})
		})

		integrationTestSharedFixture.DisposeTest()
	})

	integrationTestSharedFixture.Log.Info("TearDownSuite started")
	integrationTestSharedFixture.Bus.Stop()
	time.Sleep(1 * time.Second)
}
