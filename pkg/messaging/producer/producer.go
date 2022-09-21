package producer

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
)

type Producer interface {
	Publish(ctx context.Context, message types.IMessage, metadata core.Metadata) error
	PublishWithTopicName(ctx context.Context, message types.IMessage, metadata core.Metadata, topicOrExchangeName string) error
}
