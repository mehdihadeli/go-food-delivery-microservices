package persistmessage

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/types"

	uuid "github.com/satori/go.uuid"
)

type MessagePersistenceService interface {
	Add(ctx context.Context, storeMessage *StoreMessage) error
	Update(ctx context.Context, storeMessage *StoreMessage) error
	ChangeState(
		ctx context.Context,
		messageID uuid.UUID,
		status MessageStatus,
	) error
	GetAllActive(ctx context.Context) ([]*StoreMessage, error)
	GetByFilter(
		ctx context.Context,
		predicate func(*StoreMessage) bool,
	) ([]*StoreMessage, error)
	GetById(ctx context.Context, id uuid.UUID) (*StoreMessage, error)
	Remove(ctx context.Context, storeMessage *StoreMessage) (bool, error)
	CleanupMessages(ctx context.Context) error
	Process(messageID string, ctx context.Context) error
	ProcessAll(ctx context.Context) error
	AddPublishMessage(
		messageEnvelope types.MessageEnvelope,
		ctx context.Context,
	) error
	AddReceivedMessage(
		messageEnvelope types.MessageEnvelope,
		ctx context.Context,
	) error
	//AddInternalMessage(
	//	internalCommand InternalMessage,
	//	ctx context.Context,
	//) error
}
