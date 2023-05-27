package contracts

import "context"

type OrdersModuleConfigurator interface {
	ConfigureOrdersModule(ctx context.Context) error
}
