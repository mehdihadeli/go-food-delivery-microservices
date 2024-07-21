package postgres

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"

	"go.uber.org/fx"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("postgrespxgfx",
	fx.Provide(NewPgx, provideConfig),
	fx.Invoke(registerHooks),
)

func registerHooks(lc fx.Lifecycle, pgxClient *Pgx, logger logger.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pgxClient.Close()
			logger.Info("Pgx postgres connection closed gracefully")

			return nil
		},
	})
}
