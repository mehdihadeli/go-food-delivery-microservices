package consumer

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/test/hypothesis"
)

type RabbitMQFakeTestConsumerHandler[T any] struct {
	isHandled  bool
	hypothesis hypothesis.Hypothesis[T]
}

func NewRabbitMQFakeTestConsumerHandler[T any](hypothesis hypothesis.Hypothesis[T]) *RabbitMQFakeTestConsumerHandler[T] {
	return &RabbitMQFakeTestConsumerHandler[T]{
		hypothesis: hypothesis,
	}
}

func (f *RabbitMQFakeTestConsumerHandler[T]) Handle(ctx context.Context, consumeContext types.MessageConsumeContext) error {
	f.isHandled = true
	if f.hypothesis != nil {
		m, ok := consumeContext.Message().(T)
		if !ok {
			f.hypothesis.Test(ctx, *new(T))
		}
		f.hypothesis.Test(ctx, m)
	}

	return nil
}

func (f *RabbitMQFakeTestConsumerHandler[T]) IsHandled() bool {
	return f.isHandled
}
