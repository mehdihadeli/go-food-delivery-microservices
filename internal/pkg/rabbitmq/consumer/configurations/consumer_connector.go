//go:build go1.18

package configurations

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/consumer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
)

type RabbitMQConsumerConnector interface {
	consumer.ConsumerConnector
	ConnectRabbitMQConsumer(messageType types.IMessage, consumerBuilderFunc RabbitMQConsumerConfigurationBuilderFuc) error
}
