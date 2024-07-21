package rabbitmq

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/configurations"
	producerConfigurations "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/producer/configurations"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/events/integrationevents"
)

func ConfigProductsRabbitMQ(
	builder configurations.RabbitMQConfigurationBuilder,
) {
	builder.AddProducer(
		integrationevents.ProductCreatedV1{},
		func(builder producerConfigurations.RabbitMQProducerConfigurationBuilder) {
		},
	)
}
