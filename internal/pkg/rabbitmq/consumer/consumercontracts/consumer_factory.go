package consumercontracts

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/consumer"
	messagingTypes "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/types"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/consumer/configurations"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/rabbitmq/types"
)

type ConsumerFactory interface {
	CreateConsumer(
		consumerConfiguration *configurations.RabbitMQConsumerConfiguration,
		isConsumedNotifications ...func(message messagingTypes.IMessage),
	) (consumer.Consumer, error)

	Connection() types.IConnection
}
