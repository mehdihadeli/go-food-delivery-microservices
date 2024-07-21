package test

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/fxapp/test"

	"go.uber.org/fx/fxtest"
)

type OrdersTestApplicationBuilder struct {
	contracts.ApplicationBuilder
	tb fxtest.TB
}

func NewOrdersTestApplicationBuilder(tb fxtest.TB) *OrdersTestApplicationBuilder {
	return &OrdersTestApplicationBuilder{
		ApplicationBuilder: test.NewTestApplicationBuilder(tb),
		tb:                 tb,
	}
}

func (a *OrdersTestApplicationBuilder) Build() *OrdersTestApplication {
	return NewOrdersTestApplication(
		a.tb,
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
