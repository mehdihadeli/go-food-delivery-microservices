package contracts

import "context"

type OrdersServiceConfigurator interface {
	ConfigureOrdersService(ctx context.Context) (OrderServiceConfigurations, error)
}
