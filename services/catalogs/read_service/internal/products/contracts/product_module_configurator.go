package contracts

import "context"

type ProductsModuleConfigurator interface {
	ConfigureProductsModule(ctx context.Context) error
}
