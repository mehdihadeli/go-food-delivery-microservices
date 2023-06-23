package test

import (
	"go.uber.org/fx/fxtest"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
)

type TestApplicationBuilder struct {
	contracts.ApplicationBuilder
	TB fxtest.TB
}

func NewTestApplicationBuilder(tb fxtest.TB) *TestApplicationBuilder {
	return &TestApplicationBuilder{
		TB:                 tb,
		ApplicationBuilder: fxapp.NewApplicationBuilder(environemnt.Test),
	}
}

func (a *TestApplicationBuilder) Build() contracts.Application {
	app := NewTestApplication(
		a.TB,
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		environemnt.Test,
	)

	return app
}
