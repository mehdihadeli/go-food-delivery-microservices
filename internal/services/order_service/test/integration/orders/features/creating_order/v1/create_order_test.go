//go:build integration
// +build integration

package v1

import (
	"context"
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/hypothesis"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/dtos/v1"
	createOrderCommandV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/commands"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/dtos"
	integrationEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/creating_order/v1/events/integration_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/models/orders/read_models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/test_fixtures/integration"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var integrationFixture *integration.IntegrationTestSharedFixture

func TestCreateOrder(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewIntegrationTestSharedFixture(t)
	RunSpecs(t, "Create Order Integration Tests")
}

var _ = Describe("Create Order Feature", func() {
	var (
		ctx          context.Context
		err          error
		command      *createOrderCommandV1.CreateOrder
		result       *dtos.CreateOrderResponseDto
		createdOrder *read_models.OrderReadModel
		// id            string
		shouldPublish hypothesis.Hypothesis[*integrationEvents.OrderCreatedV1]
	)

	_ = BeforeEach(func() {
		By("Seeding the required data")
		integrationFixture.InitializeTest()

		// id = integrationFixture.Items[0].OrderId
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

	// "Scenario" for testing the creation of a new order
	Describe("Creating a new order in EventStoreDB", func() {
		BeforeEach(func() {
			command, err = createOrderCommandV1.NewCreateOrder(
				[]*dtosV1.ShopItemDto{
					{
						Quantity:    uint64(gofakeit.Number(1, 10)),
						Description: gofakeit.AdjectiveDescriptive(),
						Price:       gofakeit.Price(100, 10000),
						Title:       gofakeit.Name(),
					},
				},
				gofakeit.Email(),
				gofakeit.Address().Address,
				time.Now(),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(command).ToNot(BeNil())
		})
		When("the CreateOrder command is executed for non-existing order", func() {
			BeforeEach(func() {
				// "When" step for executing the createOrderCommand
				result, err = mediatr.Send[*createOrderCommandV1.CreateOrder, *dtos.CreateOrderResponseDto](
					ctx,
					command,
				)
			})
			// "Then" step for expected behavior
			It("Should create the order successfully", func() {
				// "Then" step for assertions
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(command.OrderId).To(Equal(result.OrderId))
			})
		})
	})

	// "Scenario" for testing the creation of a new order in MongoDB Read
	Describe("Creating a new order in MongoDB Read", func() {
		BeforeEach(func() {
			command, err = createOrderCommandV1.NewCreateOrder(
				[]*dtosV1.ShopItemDto{
					{
						Quantity:    uint64(gofakeit.Number(1, 10)),
						Description: gofakeit.AdjectiveDescriptive(),
						Price:       gofakeit.Price(100, 10000),
						Title:       gofakeit.Name(),
					},
				},
				gofakeit.Email(),
				gofakeit.Address().Address,
				time.Now(),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(command).ToNot(BeNil())
		})
		// "When" step for creating a new order
		When("the CreateOrder command is executed for non-existing order", func() {
			BeforeEach(func() {
				// "When" step for executing the createOrderCommand
				result, err = mediatr.Send[*createOrderCommandV1.CreateOrder, *dtos.CreateOrderResponseDto](
					context.Background(),
					command,
				)
			})

			It("Should create the order successfully", func() {
				// "Then" step for assertions
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
			})

			// "Then" step for expected behavior
			It("Should retrieve created order in MongoDB Read database", func() {
				// Use a utility function to wait until the order is available in MongoDB Read
				err = testUtils.WaitUntilConditionMet(func() bool {
					createdOrder, err = integrationFixture.OrderMongoRepository.GetOrderByOrderId(ctx, result.OrderId)
					Expect(err).ToNot(HaveOccurred())
					return createdOrder != nil
				})

				Expect(err).To(BeNil())
			})
		})
	})

	// "Scenario" for testing the publishing of an "OrderCreated" event
	Describe("Publishing an OrderCreated event to the message broker when order saved successfully", func() {
		BeforeEach(func() {
			shouldPublish = messaging.ShouldProduced[*integrationEvents.OrderCreatedV1](
				ctx,
				integrationFixture.Bus, nil,
			)

			command, err = createOrderCommandV1.NewCreateOrder(
				[]*dtosV1.ShopItemDto{
					{
						Quantity:    uint64(gofakeit.Number(1, 10)),
						Description: gofakeit.AdjectiveDescriptive(),
						Price:       gofakeit.Price(100, 10000),
						Title:       gofakeit.Name(),
					},
				},
				gofakeit.Email(),
				gofakeit.Address().Address,
				time.Now(),
			)

			Expect(err).ToNot(HaveOccurred())
			Expect(command).ToNot(BeNil())
		})

		// "When" step for creating and sending an order
		When("CreateOrder command is executed for non-existing order", func() {
			BeforeEach(func() {
				// "When" step for executing the createOrderCommand
				result, err = mediatr.Send[*createOrderCommandV1.CreateOrder, *dtos.CreateOrderResponseDto](
					context.Background(),
					command,
				)
			})

			It("Should return no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("Should return not nil result", func() {
				Expect(result).ToNot(BeNil())
			})

			It("Should publish OrderCreated event to the broker", func() {
				// ensuring message published to the rabbitmq broker
				shouldPublish.Validate(ctx, "there is no published message", time.Second*30)
			})
		})
	})
})
