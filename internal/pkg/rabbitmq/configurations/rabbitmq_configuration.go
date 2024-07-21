package configurations

import (
	consumerConfigurations "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/consumer/configurations"
	producerConfigurations "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/producer/configurations"
)

type RabbitMQConfiguration struct {
	ProducersConfigurations []*producerConfigurations.RabbitMQProducerConfiguration
	ConsumersConfigurations []*consumerConfigurations.RabbitMQConsumerConfiguration
}
