package gormdbcontext

import (
	"context"

	defaultlogger "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/defaultlogger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/contracts"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/helpers/gormextensions"

	"gorm.io/gorm"
)

type gormDBContext struct {
	db *gorm.DB
}

func NewGormDBContext(db *gorm.DB) contracts.GormDBContext {
	c := &gormDBContext{db: db}

	return c
}

func (c *gormDBContext) DB() *gorm.DB {
	return c.db
}

// WithTx creates a transactional DBContext with getting tx-gorm from the ctx. This will throw an error if the transaction does not exist.
func (c *gormDBContext) WithTx(
	ctx context.Context,
) (contracts.GormDBContext, error) {
	tx, err := gormextensions.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return NewGormDBContext(tx), nil
}

// WithTxIfExists creates a transactional DBContext with getting tx-gorm from the ctx. not throw an error if the transaction is not existing and returns an existing database.
func (c *gormDBContext) WithTxIfExists(
	ctx context.Context,
) contracts.GormDBContext {
	tx := gormextensions.GetTxFromContextIfExists(ctx)
	if tx == nil {
		return c
	}

	return NewGormDBContext(tx)
}

func (c *gormDBContext) RunInTx(
	ctx context.Context,
	action contracts.ActionFunc,
) error {
	// https://gorm.io/docs/transactions.html#Transaction
	tx := c.DB().WithContext(ctx).Begin()

	defaultlogger.GetLogger().Info("beginning database transaction")

	gormContext := gormextensions.SetTxToContext(ctx, tx)
	ctx = gormContext

	defer func() {
		if r := recover(); r != nil {
			tx.WithContext(ctx).Rollback()

			if err, _ := r.(error); err != nil {
				defaultlogger.GetLogger().Errorf(
					"panic tn the transaction, rolling back transaction with panic err: %+v",
					err,
				)
			} else {
				defaultlogger.GetLogger().Errorf("panic tn the transaction, rolling back transaction with panic message: %+v", r)
			}
		}
	}()

	err := action(ctx, c)
	if err != nil {
		defaultlogger.GetLogger().Error("rolling back transaction")
		tx.WithContext(ctx).Rollback()

		return err
	}

	defaultlogger.GetLogger().Info("committing transaction")

	if err = tx.WithContext(ctx).Commit().Error; err != nil {
		defaultlogger.GetLogger().Errorf("transaction commit error: %+v", err)
	}

	return err
}
