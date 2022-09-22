package consumer

import (
	"context"
)

type Consumer interface {
	Consume(ctx context.Context) error
	UnConsume(ctx context.Context) error
}
