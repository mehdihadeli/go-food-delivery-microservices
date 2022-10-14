package consumer

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
)

type RabbitMQFakeTestConsumerHandler struct {
	isHandled bool
}

func NewRabbitMQFakeTestConsumerHandler() *RabbitMQFakeTestConsumerHandler {
	return &RabbitMQFakeTestConsumerHandler{}
}

func (f *RabbitMQFakeTestConsumerHandler) Handle(ctx context.Context, consumeContext types.MessageConsumeContext) error {
	f.isHandled = true
	return nil
}

func (f *RabbitMQFakeTestConsumerHandler) IsHandled() bool {
	return f.isHandled
}
