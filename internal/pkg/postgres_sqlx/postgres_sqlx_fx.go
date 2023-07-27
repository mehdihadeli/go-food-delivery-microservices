package postgressqlx

import (
	"context"

	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("postgressqlxfx",
	fx.Provide(NewSqlxConn, provideConfig),
	fx.Invoke(registerHooks),
)

func registerHooks(lc fx.Lifecycle, pgxClient *Sqlx, logger logger.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pgxClient.Close()
			logger.Info("Sqlx postgres connection closed gracefully")

			return nil
		},
	})
}
