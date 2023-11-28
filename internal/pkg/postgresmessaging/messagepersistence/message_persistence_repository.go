package messagepersistence

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/messaging/persistmessage"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"

	"gorm.io/gorm"
)

type postgresMessagePersistenceRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewMessagePersistenceRepository(
	db *gorm.DB,
	l logger.Logger,
) persistmessage.MessagePersistenceRepository {
	return &postgresMessagePersistenceRepository{db: db, logger: l}
}

func (m *postgresMessagePersistenceRepository) Add(
	ctx context.Context,
	storeMessage *persistmessage.StoreMessage,
) error {
	result := m.db.WithContext(ctx).Create(storeMessage)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *postgresMessagePersistenceRepository) Update(
	ctx context.Context,
	storeMessage *persistmessage.StoreMessage,
) error {
	// TODO implement me
	panic("implement me")
}

func (m *postgresMessagePersistenceRepository) ChangeState(
	ctx context.Context,
	messageID string,
	status persistmessage.MessageStatus,
) error {
	// TODO implement me
	panic("implement me")
}

func (m *postgresMessagePersistenceRepository) GetAll(
	ctx context.Context,
) ([]*persistmessage.StoreMessage, error) {
	// TODO implement me
	panic("implement me")
}

func (m *postgresMessagePersistenceRepository) GetByFilter(
	ctx context.Context,
	predicate func(*persistmessage.StoreMessage) bool,
) ([]*persistmessage.StoreMessage, error) {
	// TODO implement me
	panic("implement me")
}

func (m *postgresMessagePersistenceRepository) GetById(
	ctx context.Context,
	id string,
) (*persistmessage.StoreMessage, error) {
	// TODO implement me
	panic("implement me")
}

func (m *postgresMessagePersistenceRepository) Remove(
	ctx context.Context,
	storeMessage *persistmessage.StoreMessage,
) (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (m *postgresMessagePersistenceRepository) CleanupMessages() {
	// TODO implement me
	panic("implement me")
}
