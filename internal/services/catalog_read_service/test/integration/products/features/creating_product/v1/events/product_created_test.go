//go:build integration
// +build integration

package events

// https://github.com/smartystreets/goconvey/wiki

import (
	"context"
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
	externalEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/creating_product/v1/events/integration_events/external_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/shared/test_fixture/integration"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProductCreatedConsumer(t *testing.T) {
	// Setup and initialization code here.
	integrationTestSharedFixture := integration.NewIntegrationTestSharedFixture(t)
	// in test mode we set rabbitmq `AutoStart=false` in configuration in rabbitmqOptions, so we should run rabbitmq bus manually
	integrationTestSharedFixture.Bus.Start(context.Background())
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)

	Convey("Product Created Feature", t, func() {
		// will execute with each subtest
		integrationTestSharedFixture.InitializeTest()
		ctx := context.Background()

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey("Consume ProductCreated event by consumer", func() {
			fakeProduct := &externalEvents.ProductCreatedV1{
				Message:     types.NewMessage(uuid.NewV4().String()),
				ProductId:   uuid.NewV4().String(),
				Name:        gofakeit.FirstName(),
				Price:       gofakeit.Price(150, 6000),
				CreatedAt:   time.Now(),
				Description: gofakeit.EmojiDescription(),
			}
			hypothesis := messaging.ShouldConsume[*externalEvents.ProductCreatedV1](
				ctx,
				integrationTestSharedFixture.Bus,
				nil,
			)

			Convey("When a ProductCreated event consumed", func() {
				err := integrationTestSharedFixture.Bus.PublishMessage(ctx, fakeProduct, nil)
				So(err, ShouldBeNil)

				Convey("Then it should consume the ProductCreated event", func() {
					hypothesis.Validate(ctx, "there is no consumed message", 30*time.Second)
				})
			})
		})

		Convey("Create product in mongo database when a ProductCreated event consumed", func() {
			fakeProduct := &externalEvents.ProductCreatedV1{
				Message:     types.NewMessage(uuid.NewV4().String()),
				ProductId:   uuid.NewV4().String(),
				Name:        gofakeit.FirstName(),
				Price:       gofakeit.Price(150, 6000),
				CreatedAt:   time.Now(),
				Description: gofakeit.EmojiDescription(),
			}

			Convey("When a ProductCreated event consumed", func() {
				err := integrationTestSharedFixture.Bus.PublishMessage(ctx, fakeProduct, nil)
				So(err, ShouldBeNil)

				Convey("It should store product in the mongo database", func() {
					ctx := context.Background()
					pid := uuid.NewV4().String()
					productCreated := &externalEvents.ProductCreatedV1{
						Message:     types.NewMessage(uuid.NewV4().String()),
						ProductId:   pid,
						CreatedAt:   time.Now(),
						Name:        gofakeit.Name(),
						Price:       gofakeit.Price(150, 6000),
						Description: gofakeit.AdjectiveDescriptive(),
					}

					err := integrationTestSharedFixture.Bus.PublishMessage(ctx, productCreated, nil)
					So(err, ShouldBeNil)

					var product *models.Product

					err = testUtils.WaitUntilConditionMet(func() bool {
						product, err = integrationTestSharedFixture.ProductRepository.GetProductByProductId(ctx, pid)

						return err == nil && product != nil
					})

					So(err, ShouldBeNil)
					So(product, ShouldNotBeNil)
					So(product.ProductId, ShouldEqual, pid)
				})
			})
		})

		integrationTestSharedFixture.DisposeTest()
	})

	integrationTestSharedFixture.Log.Info("TearDownSuite started")
	integrationTestSharedFixture.Bus.Stop()
	time.Sleep(1 * time.Second)
}
