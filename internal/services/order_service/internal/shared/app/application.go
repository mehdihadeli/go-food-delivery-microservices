package app

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/configurations/orders"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

type OrdersApplication struct {
	*orders.OrdersServiceConfigurator
}

func NewOrdersApplication(
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment environemnt.Environment,
) *OrdersApplication {
	app := fxapp.NewApplication(providers, decorates, options, logger, environment)
	return &OrdersApplication{
		OrdersServiceConfigurator: orders.NewOrdersServiceConfigurator(app),
	}
}
