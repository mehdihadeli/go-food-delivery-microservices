package persistmessage

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type MessagePersistenceRepository interface {
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
}
