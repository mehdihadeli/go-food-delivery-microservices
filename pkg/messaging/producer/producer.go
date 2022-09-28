package producer

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/metadata"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
)

type Producer interface {
	PublishMessage(ctx context.Context, message types.IMessage, meta metadata.Metadata) error
	PublishMessageWithTopicName(ctx context.Context, message types.IMessage, meta metadata.Metadata, topicOrExchangeName string) error
}
