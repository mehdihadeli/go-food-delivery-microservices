package messagepersistence

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgresgorm/helpers"

	"gorm.io/gorm"
)

type PostgresMessagePersistenceDBContextActionFunc func(ctx context.Context, messagePersistenceDBContext *PostgresMessagePersistenceDBContext) error

type PostgresMessagePersistenceDBContext struct {
	*gorm.DB
}

func NewPostgresMessagePersistenceDBContext(
	db *gorm.DB,
) *PostgresMessagePersistenceDBContext {
	c := &PostgresMessagePersistenceDBContext{DB: db}

	return c
}

// WithTx creates a transactional DBContext with getting tx-gorm from the ctx. This will throw an error if the transaction does not exist.
func (c *PostgresMessagePersistenceDBContext) WithTx(
	ctx context.Context,
) (*PostgresMessagePersistenceDBContext, error) {
	tx, err := helpers.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return NewPostgresMessagePersistenceDBContext(tx), nil
}

// WithTxIfExists creates a transactional DBContext with getting tx-gorm from the ctx. not throw an error if the transaction is not existing and returns an existing database.
func (c *PostgresMessagePersistenceDBContext) WithTxIfExists(
	ctx context.Context,
) *PostgresMessagePersistenceDBContext {
	tx := helpers.GetTxFromContextIfExists(ctx)
	if tx == nil {
		return c
	}

	return NewPostgresMessagePersistenceDBContext(tx)
}
