package rabbitmq

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	producerConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer/configurations"
	createProductIntegrationEvents "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/v1/events/integration_events"
)

func ConfigProductsRabbitMQ(builder configurations.RabbitMQConfigurationBuilder) {
	builder.AddProducer(
		createProductIntegrationEvents.ProductCreatedV1{},
		func(builder producerConfigurations.RabbitMQProducerConfigurationBuilder) {
		})
}
