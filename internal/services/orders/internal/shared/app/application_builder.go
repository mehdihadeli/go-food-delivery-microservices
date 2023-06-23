package app

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
)

type OrdersApplicationBuilder struct {
	*fxapp.ApplicationBuilder
}

func NewOrdersApplicationBuilder() *OrdersApplicationBuilder {
	return &OrdersApplicationBuilder{fxapp.NewApplicationBuilder()}
}

func (a *OrdersApplicationBuilder) Build() *OrdersApplication {
	return NewOrdersApplication(a.Providers, a.Options, a.Logger, a.Environment)
}
