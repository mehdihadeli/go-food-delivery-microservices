package consumer

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
)

type ConsumerConnector interface {
	// ConnectConsumerHandler adds a consumer handler to existing consumer
	ConnectConsumerHandler(messageType types.IMessage, handler ConsumerHandler)
	// ConnectConsumer creates a consumer and add handler to its handlers
	ConnectConsumer(messageType types.IMessage, handler ConsumerHandler) error
}
