package data

import "context"

type CatalogUnitOfWorkActionFunc func(catalogContext CatalogContext) error

type CatalogUnitOfWork interface {
	// Do execute the given CatalogUnitOfWorkActionFunc atomically (inside a DB transaction).
	Do(ctx context.Context, action CatalogUnitOfWorkActionFunc) error
}
