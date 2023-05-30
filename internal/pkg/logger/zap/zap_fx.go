package zap

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

// Module provided to fx
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("zapfx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		NewZapLogger,
		fx.Annotate(
			NewZapLogger,
			fx.As(new(logger.Logger))),
		provideLogConfig,
	),
	fx.WithLogger(func(log ZapLogger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: log.InternalLogger()}
	},
	))

func provideLogConfig() (*logger.LogConfig, error) {
	return config.BindConfigKey[*logger.LogConfig]("logger")
}
