package producer

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/messaging/producer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	producerConfigurations "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/producer/configurations"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/producer/producercontracts"
	types2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/types"
)

type producerFactory struct {
	connection      types2.IConnection
	logger          logger.Logger
	eventSerializer serializer.EventSerializer
	rabbitmqOptions *config.RabbitmqOptions
}

func NewProducerFactory(
	rabbitmqOptions *config.RabbitmqOptions,
	connection types2.IConnection,
	eventSerializer serializer.EventSerializer,
	l logger.Logger,
) producercontracts.ProducerFactory {
	return &producerFactory{
		rabbitmqOptions: rabbitmqOptions,
		logger:          l,
		connection:      connection,
		eventSerializer: eventSerializer,
	}
}

func (p *producerFactory) CreateProducer(
	rabbitmqProducersConfiguration map[string]*producerConfigurations.RabbitMQProducerConfiguration,
	isProducedNotifications ...func(message types.IMessage),
) (producer.Producer, error) {
	return NewRabbitMQProducer(
		p.rabbitmqOptions,
		p.connection,
		rabbitmqProducersConfiguration,
		p.logger,
		p.eventSerializer,
		isProducedNotifications...)
}
