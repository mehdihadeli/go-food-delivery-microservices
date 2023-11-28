package persistmessage

import "context"

type MessagePersistenceRepository interface {
	Add(ctx context.Context, storeMessage *StoreMessage) error
	Update(ctx context.Context, storeMessage *StoreMessage) error
	ChangeState(
		ctx context.Context,
		messageID string,
		status MessageStatus,
	) error
	GetAll(ctx context.Context) ([]*StoreMessage, error)
	GetByFilter(
		ctx context.Context,
		predicate func(*StoreMessage) bool,
	) ([]*StoreMessage, error)
	GetById(ctx context.Context, id string) (*StoreMessage, error)
	Remove(ctx context.Context, storeMessage *StoreMessage) (bool, error)
	CleanupMessages()
}
