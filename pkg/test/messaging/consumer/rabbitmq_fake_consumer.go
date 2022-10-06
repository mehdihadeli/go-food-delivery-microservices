package consumer

import (
	"context"
	consumer2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
)

type RabbitMQFakeTestConsumer struct {
	isHandled bool
	consumer2.Consumer
}

func NewRabbitMQFakeTestConsumer() *RabbitMQFakeTestConsumer {
	return &RabbitMQFakeTestConsumer{}
}

func (f *RabbitMQFakeTestConsumer) Handle(ctx context.Context, consumeContext types.MessageConsumeContext) error {
	f.isHandled = true
	return nil
}

func (f *RabbitMQFakeTestConsumer) IsHandled() bool {
	return f.isHandled
}
