package contracts

import "context"

type CatalogsServiceConfigurator interface {
	ConfigureCatalogsService(ctx context.Context) (CatalogServiceConfigurations, error)
}
