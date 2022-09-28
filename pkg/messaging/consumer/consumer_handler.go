package consumer

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
)

type ConsumerHandlerT[T types.IMessage] interface {
	Handle(ctx context.Context, consumeContext types.MessageConsumeContextT[T]) error
}

type ConsumerHandler interface {
	Handle(ctx context.Context, consumeContext types.MessageConsumeContext) error
}
