//go:build integration
// +build integration

package v1

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/messaging/consumer"
	integrationEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/events/integration_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/events/integration_events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/test_fixtures/integration"
)

type productFeaturesIntegrationTestsv1 struct {
	*integration.IntegrationTestSharedFixture
}

func TestProductFeaturesIntegration(t *testing.T) {
	suite.Run(
		t,
		&productFeaturesIntegrationTestsv1{
			IntegrationTestSharedFixture: integration.NewIntegrationTestSharedFixture(t),
		},
	)
}

func (p *productFeaturesIntegrationTestsv1) BeforeTest(suiteName, testName string) {
	if testName == "Test_Should_Consume_Product_Created_With_New_Consumer_From_Broker" {
		p.Bus.Stop()
	}
}

func (p *productFeaturesIntegrationTestsv1) SetupSuite() {
	// we don't have a consumer in this service, so we simulate one consumer, register one consumer for `ProductCreatedV1` message before executing the tests
	testConsumer := consumer.NewRabbitMQFakeTestConsumerHandler[*integrationEvents.ProductCreatedV1]()
	err := p.Bus.ConnectConsumerHandler(&integrationEvents.ProductCreatedV1{}, testConsumer)
	p.Require().NoError(err)

	// we don't have a consumer in this service, so we simulate one consumer, register one consumer for `ProductUpdatedV1` message before executing the tests
	productUpdatedConsumer := consumer.NewRabbitMQFakeTestConsumerHandler[*integration_events.ProductUpdatedV1]()
	err = p.Bus.ConnectConsumerHandler(
		&integration_events.ProductUpdatedV1{},
		productUpdatedConsumer,
	)
	p.Require().NoError(err)

	// in test mode we set rabbitmq `AutoStart=false`, so we should run rabbitmq bus manually
	p.Bus.Start(context.Background())
	// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
	time.Sleep(1 * time.Second)
}

func (p *productFeaturesIntegrationTestsv1) TearDownSuite() {
	p.Bus.Stop()
}
