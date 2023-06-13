package fxapp

import (
	"go.uber.org/fx"
)

type Application struct {
	providers []interface{}
	invokes   []interface{}
	options   []fx.Option
}

func NewApplication(providers []interface{}, options []fx.Option) *Application {
	return &Application{providers: providers, options: options}
}

func (a *Application) ResolveFunc(function interface{}) {
	a.invokes = append(a.invokes, function)
}

func (a *Application) RegisterHook(function interface{}) {
	a.invokes = append(a.invokes, function)
}

func (a *Application) Run() {
	// build phase of container will do in this stage, containing provides and invokes but app not started yet and will be started in the future with `fxApp.Run`
	fxApp := CreateFxApp(a.providers, a.invokes, a.options...)

	// running phase will do in this stage and all register event hooks like OnStart and OnStop
	fxApp.Run()
}
