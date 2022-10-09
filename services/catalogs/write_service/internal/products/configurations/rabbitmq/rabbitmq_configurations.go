package rabbitmq

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/configurations"
	producerConfigurations "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/producer/configurations"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/events/integration/v1"
)

func ConfigProductsRabbitMQ(builder configurations.RabbitMQConfigurationBuilder) {
	builder.AddProducer(
		v1.ProductCreatedV1{},
		func(builder producerConfigurations.RabbitMQProducerConfigurationBuilder) {
		})
}
