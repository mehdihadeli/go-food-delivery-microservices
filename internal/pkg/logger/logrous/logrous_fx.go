package logrous

import (
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

// Module provided to fx
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("logrousfx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		fx.Annotate(
			NewLogrusLogger,
			fx.As(new(logger.Logger)),
		),

		provideLogConfig,
	))

func provideLogConfig() (*logger.LogConfig, error) {
	return config.BindConfigKey[*logger.LogConfig]("logger")
}
