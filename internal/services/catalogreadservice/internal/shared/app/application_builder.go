package app

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/fxapp"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/fxapp/contracts"
)

type CatalogsReadApplicationBuilder struct {
	contracts.ApplicationBuilder
}

func NewCatalogsReadApplicationBuilder() *CatalogsReadApplicationBuilder {
	builder := &CatalogsReadApplicationBuilder{fxapp.NewApplicationBuilder()}

	return builder
}

func (a *CatalogsReadApplicationBuilder) Build() *CatalogsReadApplication {
	return NewCatalogsReadApplication(
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
