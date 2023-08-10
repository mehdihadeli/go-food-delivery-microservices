package test

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/test"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/app"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/shared/configurations/orders"
)

type OrdersTestApplication struct {
	*app.OrdersApplication
	tb fxtest.TB
}

func NewOrdersTestApplication(
	tb fxtest.TB,
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment environemnt.Environment,
) *OrdersTestApplication {
	testApp := test.NewTestApplication(
		tb,
		providers,
		decorates,
		options,
		logger,
		environment,
	)

	orderApplication := &app.OrdersApplication{
		OrdersServiceConfigurator: orders.NewOrdersServiceConfigurator(testApp),
	}

	return &OrdersTestApplication{
		OrdersApplication: orderApplication,
		tb:                tb,
	}
}
