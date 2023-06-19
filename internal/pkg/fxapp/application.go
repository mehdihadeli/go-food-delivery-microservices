package fxapp

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

type Application struct {
	providers []interface{}
	invokes   []interface{}
	options   []fx.Option
	Logger    logger.Logger
}

func NewApplication(
	providers []interface{},
	options []fx.Option,
	logger logger.Logger,
) *Application {
	return &Application{providers: providers, options: options, Logger: logger}
}

func (a *Application) ResolveFunc(function interface{}) {
	a.invokes = append(a.invokes, function)
}

func (a *Application) ResolveFuncWithParamTag(function interface{}, paramTagName string) {
	a.invokes = append(a.invokes, fx.Annotate(function, fx.ParamTags(paramTagName)))
}

func (a *Application) RegisterHook(function interface{}) {
	a.invokes = append(a.invokes, function)
}

func (a *Application) Run() {
	// build phase of container will do in this stage, containing provides and invokes but app not started yet and will be started in the future with `fxApp.Run`
	fxApp := CreateFxApp(a.Logger, a.providers, a.invokes, a.options...)

	// running phase will do in this stage and all register event hooks like OnStart and OnStop
	// instead of run for handling start and stop and create a ctx and cancel we can handle them manually with fx.start and fx.stop
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
