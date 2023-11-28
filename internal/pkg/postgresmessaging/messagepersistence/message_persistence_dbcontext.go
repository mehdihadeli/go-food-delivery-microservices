package messagepersistence

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgresGorm/helpers"

	"gorm.io/gorm"
)

type PostgresMessagePersistenceDBContextActionFunc func(ctx context.Context, messagePersistenceDBContext *PostgresMessagePersistenceDBContext) error

type PostgresMessagePersistenceDBContext struct {
	*gorm.DB
	logger logger.Logger
}

func NewPostgresMessagePersistenceDBContext(
	db *gorm.DB,
	log logger.Logger,
) *PostgresMessagePersistenceDBContext {
	c := &PostgresMessagePersistenceDBContext{DB: db, logger: log}

	return c
}

// WithTx creates a transactional DBContext with getting tx-gorm from the ctx
func (c *PostgresMessagePersistenceDBContext) WithTx(
	ctx context.Context,
) (*PostgresMessagePersistenceDBContext, error) {
	tx, err := helpers.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return NewPostgresMessagePersistenceDBContext(tx, c.logger), nil
}
