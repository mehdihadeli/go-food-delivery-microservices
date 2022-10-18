package uow

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/data"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"gorm.io/gorm"
)

type gormUnitOfWork struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewGormUnitOfWork(db *gorm.DB, logger logger.Logger) data.UnitOfWork {
	fmt.Println(&db)
	return &gormUnitOfWork{db: db, logger: logger}
}

func (g *gormUnitOfWork) SaveWithTx(ctx context.Context, action data.UnitOfWorkActionFunc) error {
	var err error
	preserveGormInstance := *g.db

	// using ctx for support cancellation inner transaction
	tx := g.db.WithContext(ctx).Begin()
	g.logger.Info("beginning database transaction")

	// change all repository gorm referenced instances inner transaction
	*g.db = *tx

	defer func() {
		r := recover()
		if r != nil {
			tx.WithContext(ctx).Rollback()
			err, _ = r.(error)
			if err != nil {
				g.logger.Errorf("panic tn the transaction, rolling back transaction with panic err: %+v", err)
			} else {
				g.logger.Errorf("panic tn the transaction, rolling back transaction with panic message: %+v", r)
			}
		}

		// back to gorm instance before starting transaction
		*g.db = preserveGormInstance
	}()
	err = action()
	if err == nil {
		if txErr := tx.WithContext(ctx).Commit().Error; txErr != nil {
			g.logger.Errorf("transaction commit error: %+v", txErr)
			err = errors.WrapIf(err, txErr.Error())
			return err
		}
		g.logger.Info("transaction commit succeeded")
		return nil
	} else {
		g.logger.Errorf("rolling back transaction, err: %+v", err)
		tx.WithContext(ctx).Rollback()
	}

	return err
}
