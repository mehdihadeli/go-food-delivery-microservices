package fxapp

import (
	"go.uber.org/fx"
)

type ApplicationBuilder struct {
	Providers []interface{}
	Options   []fx.Option
}

func NewApplicationBuilder() *ApplicationBuilder {
	return &ApplicationBuilder{}
}

func (a *ApplicationBuilder) ProvideModule(module fx.Option) {
	a.Options = append(a.Options, module)
}

func (a *ApplicationBuilder) ProvideFunc(constructors ...interface{}) {
	a.Providers = append(a.Providers, constructors...)
}

func (a *ApplicationBuilder) Build() *Application {
	app := NewApplication(a.Providers, a.Options)

	return app
}
