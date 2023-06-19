package fxapp

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	config2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/external/fxlog"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/logrous"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/models"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/zap"
)

func CreateFxApp(
	logger logger.Logger,
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

	var logModule fx.Option
	logoption, err := config2.ProvideLogConfig()
	if err != nil || logoption == nil {
		logModule = zap.ModuleFunc(logger)
	} else if logoption.LogType == models.Logrus {
		logModule = logrous.ModuleFunc(logger)
	} else {
		logModule = zap.ModuleFunc(logger)
	}

	// build phase of container will do in this stage, containing provides and invokes but app not started yet and will be started in the future with `fxApp.Run`
	fxApp := fx.New( // setup fxlog logger
		config.Module,
		logModule,
		fxlog.FxLogger,
		fx.ErrorHook(NewFxErrorHandler(logger)),
		AppModule,
	)

	return fxApp
}
