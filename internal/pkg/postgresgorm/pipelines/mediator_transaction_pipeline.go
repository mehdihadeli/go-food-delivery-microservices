package pipelines

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/cqrs"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/helpers/gormextensions"
	typeMapper "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/reflection/typemapper"

	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

type mediatorTransactionPipeline struct {
	logger logger.Logger
	db     *gorm.DB
}

func NewMediatorTransactionPipeline(
	l logger.Logger,
	db *gorm.DB,
) mediatr.PipelineBehavior {
	return &mediatorTransactionPipeline{
		logger: l,
		db:     db,
	}
}

func (m *mediatorTransactionPipeline) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {
	requestName := typeMapper.GetSnakeTypeName(request)

	_, ok := request.(cqrs.TxRequest)
	if !ok {
		return next(ctx)
	}

	var result interface{}

	// https://gorm.io/docs/transactions.html#Transaction
	tx := m.db.WithContext(ctx).Begin()

	m.logger.Infof(
		"beginning database transaction for request `%s`",
		requestName,
	)

	gormContext := gormextensions.SetTxToContext(ctx, tx)
	ctx = gormContext

	defer func() {
		if r := recover(); r != nil {
			tx.WithContext(ctx).Rollback()

			if err, _ := r.(error); err != nil {
				m.logger.Errorf(
					"panic tn the transaction, rolling back transaction with panic err: %+v",
					err,
				)
			} else {
				m.logger.Errorf("panic tn the transaction, rolling back transaction with panic message: %+v", r)
			}
		}
	}()

	middlewareResponse, err := next(ctx)
	result = middlewareResponse

	if err != nil {
		m.logger.Errorf(
			"rolling back transaction for request `%s`",
			requestName,
		)
		tx.WithContext(ctx).Rollback()

		return nil, err
	}

	m.logger.Infof("committing transaction for request `%s`", requestName)

	if err = tx.WithContext(ctx).Commit().Error; err != nil {
		m.logger.Errorf("transaction commit error: ", err)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}
