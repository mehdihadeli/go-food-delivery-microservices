package consumer

import (
	"context"
)

type Consumer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	ConnectHandler(handler ConsumerHandler)
}
