package bus

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/producer"
)

type Bus interface {
	producer.Producer
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
