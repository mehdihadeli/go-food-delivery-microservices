package app

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
)

type CatalogsWriteApplicationBuilder struct {
	*fxapp.ApplicationBuilder
}

func NewCatalogsWriteApplicationBuilder() *CatalogsWriteApplicationBuilder {
	return &CatalogsWriteApplicationBuilder{fxapp.NewApplicationBuilder()}
}

func (a *CatalogsWriteApplicationBuilder) Build() *CatalogsWriteApplication {
	return NewCatalogsWriteApplication(a.Providers, a.Options)
}
