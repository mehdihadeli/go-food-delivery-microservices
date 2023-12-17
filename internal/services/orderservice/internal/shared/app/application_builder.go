package app

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
)

type OrdersApplicationBuilder struct {
	contracts.ApplicationBuilder
}

func NewOrdersApplicationBuilder() *OrdersApplicationBuilder {
	return &OrdersApplicationBuilder{fxapp.NewApplicationBuilder()}
}

func (a *OrdersApplicationBuilder) Build() *OrdersApplication {
	return NewOrdersApplication(
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
