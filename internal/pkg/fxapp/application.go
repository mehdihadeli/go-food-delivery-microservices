package fxapp

import (
	"context"

	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/fxapp/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

type application struct {
	provides    []interface{}
	decorates   []interface{}
	invokes     []interface{}
	options     []fx.Option
	logger      logger.Logger
	fxapp       *fx.App
	environment environemnt.Environment
}

func NewApplication(
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	env environemnt.Environment,
) contracts.Application {
	return &application{
		provides:    providers,
		decorates:   decorates,
		options:     options,
		logger:      logger,
		environment: env,
	}
}

func (a *application) ResolveFunc(function interface{}) {
	a.invokes = append(a.invokes, function)
}

func (a *application) ResolveFuncWithParamTag(function interface{}, paramTagName string) {
	a.invokes = append(a.invokes, fx.Annotate(function, fx.ParamTags(paramTagName)))
}

func (a *application) RegisterHook(function interface{}) {
	a.invokes = append(a.invokes, function)
}

func (a *application) Run() {
	// build phase of container will do in this stage, containing provides and invokes but app not started yet and will be started in the future with `fxApp.Run`
	fxApp := CreateFxApp(a)

	a.fxapp = fxApp

	// running phase will do in this stage and all register event hooks like OnStart and OnStop
	// instead of run for handling start and stop and create a ctx and cancel we can handle them manually with appconfigfx.start and appconfigfx.stop
	// https://github.com/uber-go/fx/blob/v1.20.0/app.go#L573
	fxApp.Run()

	//// startup ctx just for setup dependencies about 15 seconds
	//const defaultTimeout = 15 * time.Second
	//startCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	//defer cancel()
	//
	//if err := fxApp.Start(startCtx); err != nil {
	//	os.Exit(1)
	//}
	//// wait until get a os signal
	//sig := <-fxApp.Wait()
	//exitCode := sig.ExitCode
	//// shutdown ctx just for shut down process and about 15 seconds
	//stopCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	//defer cancel()
	//
	//if err := fxApp.Stop(stopCtx); err != nil {
	//	exitCode = 1
	//}
	//
	//if exitCode != 0 {
	//	os.Exit(exitCode)
	//}
}

func (a *application) Start(ctx context.Context) error {
	// build phase of container will do in this stage, containing provides and invokes but app not started yet and will be started in the future with `fxApp.Run`
	fxApp := CreateFxApp(a)
	a.fxapp = fxApp

	return fxApp.Start(ctx)
}

func (a *application) Stop(ctx context.Context) error {
	if a.fxapp == nil {
		a.logger.Fatal("Failed to stop because application not started.")
	}
	return a.fxapp.Stop(ctx)
}

func (a *application) Wait() <-chan fx.ShutdownSignal {
	if a.fxapp == nil {
		a.logger.Fatal("Failed to wait because application not started.")
	}
	return a.fxapp.Wait()
}

func (a *application) Logger() logger.Logger {
	return a.logger
}

func (a *application) Environment() environemnt.Environment {
	return a.environment
}
