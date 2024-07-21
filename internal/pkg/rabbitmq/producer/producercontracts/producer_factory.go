package producercontracts

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/producer"
	types2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/types"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/producer/configurations"
)

type ProducerFactory interface {
	CreateProducer(
		rabbitmqProducersConfiguration map[string]*configurations.RabbitMQProducerConfiguration,
		isProducedNotifications ...func(message types2.IMessage),
	) (producer.Producer, error)
}
