package test

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/fxapp/test"

	"go.uber.org/fx/fxtest"
)

type CatalogsReadTestApplicationBuilder struct {
	contracts.ApplicationBuilder
	tb fxtest.TB
}

func NewCatalogsReadTestApplicationBuilder(tb fxtest.TB) *CatalogsReadTestApplicationBuilder {
	return &CatalogsReadTestApplicationBuilder{
		ApplicationBuilder: test.NewTestApplicationBuilder(tb),
		tb:                 tb,
	}
}

func (a *CatalogsReadTestApplicationBuilder) Build() *CatalogsReadTestApplication {
	return NewCatalogsReadTestApplication(
		a.tb,
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}
