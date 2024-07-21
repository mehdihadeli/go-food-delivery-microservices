package gormextensions

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/constants"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/contracts"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/scopes"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"

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

func GetTxFromContextIfExists(ctx context.Context) *gorm.DB {
	gCtx, gCtxOk := ctx.(*contracts.GormContext)
	if gCtxOk {
		return gCtx.Tx
	}

	tx, ok := ctx.Value(constants.TxKey).(*gorm.DB)
	if !ok {
		return nil
	}

	return tx
}

func SetTxToContext(ctx context.Context, tx *gorm.DB) *contracts.GormContext {
	newCtx := context.WithValue(ctx, constants.TxKey, tx)
	gormContext := &contracts.GormContext{Tx: tx, Context: newCtx}
	ctx = gormContext

	return gormContext
}

// Ref: https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f

func Paginate[TDataModel any, TEntity any](
	ctx context.Context,
	listQuery *utils.ListQuery,
	db *gorm.DB,
) (*utils.ListResult[TEntity], error) {
	var (
		items     []TEntity
		totalRows int64
	)

	// https://gorm.io/docs/advanced_query.html#Smart-Select-Fields
	if err := db.Scopes(scopes.FilterPaginate[TDataModel](ctx, listQuery)).Find(&items).Error; err != nil {
		return nil, errors.WrapIf(err, "error in finding products.")
	}

	return utils.NewListResult[TEntity](
		items,
		listQuery.GetSize(),
		listQuery.GetPage(),
		totalRows,
	), nil
}
