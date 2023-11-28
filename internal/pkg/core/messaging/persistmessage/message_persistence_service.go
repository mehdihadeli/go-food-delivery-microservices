package persistmessage

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/messaging/types"
)

type IMessagePersistenceService interface {
	GetByFilter(
		predicate func(StoreMessage) bool,
		ctx context.Context,
	) ([]StoreMessage, error)
	AddPublishMessage(
		messageEnvelope types.MessageEnvelopeTMessage,
		ctx context.Context,
	) error
	AddReceivedMessage(
		messageEnvelope types.MessageEnvelope,
		ctx context.Context,
	) error
	Process(messageID string, ctx context.Context) error
	ProcessAll(ctx context.Context) error
}
