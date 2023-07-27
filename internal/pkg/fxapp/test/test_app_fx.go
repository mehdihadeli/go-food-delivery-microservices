package test

import (
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/zap"
)

func CreateFxTestApp(
	tb fxtest.TB,
	provides []interface{},
	decorates []interface{},
	invokes []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment environemnt.Environment,
) *fxtest.App {
	var opts []fx.Option

	opts = append(opts, fx.Provide(provides...))

	opts = append(opts, fx.Invoke(invokes...))

	opts = append(opts, fx.Decorate(decorates...))

	options = append(options, opts...)

	AppModule := fx.Module("fxtestapp",
		options...,
	)

	duration := 60 * time.Second

	// build phase of container will do in this stage, containing provides and invokes but app not started yet and will be started in the future with `fxApp.Run`
	fxApp := fxtest.New(
		tb,
		fx.StartTimeout(duration),
		config.ModuleFunc(environment),
		zap.ModuleFunc(logger),
		AppModule,

		// fx.Decorate(rabbitmq.RabbitmqContainerDecorator(tb.(*testing.T), context.Background())),
	)

	return fxApp
}
