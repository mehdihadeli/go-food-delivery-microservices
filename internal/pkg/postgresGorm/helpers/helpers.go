package helpers

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgresGorm/constants"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgresGorm/contracts"

	"emperror.dev/errors"
	"gorm.io/gorm"
)

func GetTxFromContext(ctx context.Context) (*gorm.DB, error) {
	gCtx, gCtxOk := ctx.(*contracts.GormContext)
	if gCtxOk {
		return gCtx.Tx, nil
	}

	tx, ok := ctx.Value(constants.TxKey).(*gorm.DB)
	if !ok {
		return nil, errors.New("Transaction not found in context")
	}

	return tx, nil
}

func SetTxToContext(ctx context.Context, tx *gorm.DB) *contracts.GormContext {
	ctx = context.WithValue(ctx, constants.TxKey, tx)
	gormContext := &contracts.GormContext{Tx: tx, Context: ctx}

	return gormContext
}
