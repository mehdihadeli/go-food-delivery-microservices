package zap

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/config"
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
