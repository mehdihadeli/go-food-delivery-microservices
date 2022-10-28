package consumer

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
)

type Consumer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	ConnectHandler(handler ConsumerHandler)
	AddMessageConsumedHandler(func(message types.IMessage))
}
