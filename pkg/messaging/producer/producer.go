package producer

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
)

type Producer interface {
	Publish(ctx context.Context, topicOrExchangeName string, message types.IMessage, metadata core.Metadata) error
}
