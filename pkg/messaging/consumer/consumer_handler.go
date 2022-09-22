package consumer

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
)

type ConsumerHandler[T types.IMessage] interface {
	Handle(ctx context.Context, consumeContext types.IMessageConsumeContext[T]) error
}
