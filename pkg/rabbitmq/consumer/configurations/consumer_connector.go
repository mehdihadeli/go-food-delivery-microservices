package configurations

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
)

type RabbitMQConsumerConnector interface {
	consumer.ConsumerConnector
	ConnectRabbitMQConsumer(messageType types.IMessage, consumerBuilderFunc RabbitMQConsumerConfigurationBuilderFuc) error
}
