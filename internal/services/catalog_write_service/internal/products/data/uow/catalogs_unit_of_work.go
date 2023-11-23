package uow

// https://blog.devgenius.io/go-golang-unit-of-work-and-generics-5e9fb00ec996
// https://learn.microsoft.com/en-us/aspnet/mvc/overview/older-versions/getting-started-with-ef-5-using-mvc-4/implementing-the-repository-and-unit-of-work-patterns-in-an-asp-net-mvc-application
// https://dev.to/techschoolguru/a-clean-way-to-implement-database-transaction-in-golang-2ba

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/data/repositories"

	"gorm.io/gorm"
)

type catalogUnitOfWork[TContext contracts.CatalogContext] struct {
	logger logger.Logger
	db     *gorm.DB
	tracer tracing.AppTracer
}

func NewCatalogsUnitOfWork(
	logger logger.Logger,
	db *gorm.DB,
	tracer tracing.AppTracer,
) contracts.CatalogUnitOfWork {
	return &catalogUnitOfWork[contracts.CatalogContext]{logger: logger, db: db, tracer: tracer}
}

func (c *catalogUnitOfWork[TContext]) Do(
	ctx context.Context,
	action contracts.CatalogUnitOfWorkActionFunc,
) error {
	// https://gorm.io/docs/transactions.html#Transaction
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		catalog := &catalogContext{
			productRepository: repositories.NewPostgresProductRepository(c.logger, tx, c.tracer),
		}

		defer func() {
			r := recover()
			if r != nil {
				tx.WithContext(ctx).Rollback()
				err, _ := r.(error)
				if err != nil {
					c.logger.Errorf(
						"panic tn the transaction, rolling back transaction with panic err: %+v",
						err,
					)
				} else {
					c.logger.Errorf("panic tn the transaction, rolling back transaction with panic message: %+v", r)
				}
			}
		}()

		return action(catalog)
	})
}
