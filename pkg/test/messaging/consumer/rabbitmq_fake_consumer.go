package consumer

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	consumer2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/consumer/options"
	types2 "github.com/mehdihadeli/store-golang-microservice-sample/pkg/rabbitmq/types"
)

type RabbitMQFakeTestConsumer struct {
	isHandled bool
	consumer2.Consumer
}

func NewRabbitMQFakeTestConsumer(eventSerializer serializer.EventSerializer, logger logger.Logger, connection types2.IConnection, builderFunc func(builder *options.RabbitMQConsumerOptionsBuilder)) *RabbitMQFakeTestConsumer {
	fakeConsumer := &RabbitMQFakeTestConsumer{}

	t, err := consumer.NewRabbitMQConsumer(eventSerializer, logger, connection, builderFunc, fakeConsumer)
	if err != nil {
		return nil
	}
	fakeConsumer.Consumer = t
	return fakeConsumer
}

func (f *RabbitMQFakeTestConsumer) Handle(ctx context.Context, consumeContext types.MessageConsumeContext) error {
	f.isHandled = true
	return nil
}

func (f *RabbitMQFakeTestConsumer) IsHandled() bool {
	return f.isHandled
}
