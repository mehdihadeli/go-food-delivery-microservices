package empty

import (
	logger2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/config"

	"go.uber.org/fx"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("emptyfx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		fx.Annotate(
			EmptyLogger,
			fx.As(new(logger2.Logger)),
		),
		config.ProvideLogConfig,
	))
