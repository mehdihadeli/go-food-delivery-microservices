package consumer

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/consumer"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/types"
	serializer "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/serializer"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/config"
	consumerConfigurations "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/consumer/configurations"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/consumer/consumercontracts"
	types2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/types"
)

type consumerFactory struct {
	connection      types2.IConnection
	eventSerializer serializer.MessageSerializer
	logger          logger.Logger
	rabbitmqOptions *config.RabbitmqOptions
}

func NewConsumerFactory(
	rabbitmqOptions *config.RabbitmqOptions,
	connection types2.IConnection,
	eventSerializer serializer.MessageSerializer,
	l logger.Logger,
) consumercontracts.ConsumerFactory {
	return &consumerFactory{
		rabbitmqOptions: rabbitmqOptions,
		logger:          l,
		eventSerializer: eventSerializer,
		connection:      connection,
	}
}

func (c *consumerFactory) CreateConsumer(
	consumerConfiguration *consumerConfigurations.RabbitMQConsumerConfiguration,
	isConsumedNotifications ...func(message types.IMessage),
) (consumer.Consumer, error) {
	return NewRabbitMQConsumer(
		c.rabbitmqOptions,
		c.connection,
		consumerConfiguration,
		c.eventSerializer,
		c.logger,
		isConsumedNotifications...)
}

func (c *consumerFactory) Connection() types2.IConnection {
	return c.connection
}
