package zap

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/config"

	"go.uber.org/fx"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("zapfx",

	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		config.ProvideLogConfig,
		NewZapLogger,
		fx.Annotate(
			NewZapLogger,
			fx.As(new(logger.Logger))),
	),
)

var ModuleFunc = func(l logger.Logger) fx.Option {
	return fx.Module(
		"zapfx",

		fx.Provide(config.ProvideLogConfig),
		fx.Supply(fx.Annotate(l, fx.As(new(logger.Logger)))),
		fx.Supply(fx.Annotate(l, fx.As(new(ZapLogger)))),
	)
}
