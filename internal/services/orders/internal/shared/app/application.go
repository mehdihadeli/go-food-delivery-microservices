package app

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/shared/configurations/orders"
)

type OrdersApplication struct {
	*orders.OrdersServiceConfigurator
}

func NewOrdersApplication(
	providers []interface{},
	options []fx.Option,
	logger logger.Logger,
) *OrdersApplication {
	app := fxapp.NewApplication(providers, options, logger)
	return &OrdersApplication{
		OrdersServiceConfigurator: orders.NewOrdersServiceConfigurator(app),
	}
}
