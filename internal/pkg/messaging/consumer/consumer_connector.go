package consumer

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
)

type ConsumerConnector interface {
	// ConnectConsumerHandler adds a consumer handler to existing consumer and create a new consumer if it doesn't already exist
	ConnectConsumerHandler(messageType types.IMessage, consumerHandler ConsumerHandler) error
	// ConnectConsumer add consumer to consumers list
	ConnectConsumer(messageType types.IMessage, consumer Consumer) error
}
