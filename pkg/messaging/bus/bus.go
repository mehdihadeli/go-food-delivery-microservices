package bus

import (
	"context"
)

type Bus interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
