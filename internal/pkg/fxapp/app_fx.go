package fxapp

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/external/fxlog"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/zap"
)

func CreateFxApp(
	providers []interface{},
	invokes []interface{},
	options ...fx.Option,
) *fx.App {
	var opts []fx.Option

	opts = append(opts, fx.Provide(providers...))

	opts = append(opts, fx.Invoke(invokes...))

	options = append(options, opts...)

	AppModule := fx.Module("appfx",
		options...,
	)

	// build phase of container will do in this stage, containing provides and invokes but app not started yet and will be started in the future with `fxApp.Run`
	fxApp := fx.New( // setup fxlog logger
		zap.Module,
		config.Module,
		fxlog.FxLogger,
		fx.ErrorHook(NewFxErrorHandler()),

		AppModule,
	)
	return fxApp
}
