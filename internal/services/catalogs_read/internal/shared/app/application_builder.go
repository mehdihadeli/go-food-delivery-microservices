package app

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
)

type CatalogsReadApplicationBuilder struct {
	*fxapp.ApplicationBuilder
}

func NewCatalogsReadApplicationBuilder() *CatalogsReadApplicationBuilder {
	return &CatalogsReadApplicationBuilder{fxapp.NewApplicationBuilder()}
}

func (a *CatalogsReadApplicationBuilder) Build() *CatalogsReadApplication {
	return NewCatalogsReadApplication(a.Providers, a.Options)
}
