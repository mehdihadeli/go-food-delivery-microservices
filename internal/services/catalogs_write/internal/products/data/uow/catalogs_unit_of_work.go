package uow

// https://blog.devgenius.io/go-golang-unit-of-work-and-generics-5e9fb00ec996
// https://learn.microsoft.com/en-us/aspnet/mvc/overview/older-versions/getting-started-with-ef-5-using-mvc-4/implementing-the-repository-and-unit-of-work-patterns-in-an-asp-net-mvc-application
// https://dev.to/techschoolguru/a-clean-way-to-implement-database-transaction-in-golang-2ba

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/data/repositories"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"gorm.io/gorm"

	data2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogs/write_service/internal/products/contracts/data"
)

type catalogUnitOfWork[TContext data2.CatalogContext] struct {
	logger logger.Logger
	db     *gorm.DB
}

func (c *catalogUnitOfWork[TContext]) Do(ctx context.Context, action data2.CatalogUnitOfWorkActionFunc) error {
	// https://gorm.io/docs/transactions.html#Transaction
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		catalog := &catalogContext{
			productRepository: repositories.NewPostgresProductRepository(c.logger, tx),
		}

		defer func() {
			r := recover()
			if r != nil {
				tx.WithContext(ctx).Rollback()
				err, _ := r.(error)
				if err != nil {
					c.logger.Errorf("panic tn the transaction, rolling back transaction with panic err: %+v", err)
				} else {
					c.logger.Errorf("panic tn the transaction, rolling back transaction with panic message: %+v", r)
				}
			}
		}()

		return action(catalog)
	})
}

func NewCatalogsUnitOfWork(logger logger.Logger, db *gorm.DB) data2.CatalogUnitOfWork {
	return &catalogUnitOfWork[data2.CatalogContext]{logger: logger, db: db}
}
