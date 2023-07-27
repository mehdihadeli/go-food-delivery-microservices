package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("redisfx",
	fx.Provide(NewUniversalRedisClient, provideConfig),
	fx.Invoke(registerHooks),
)

func registerHooks(lc fx.Lifecycle, client redis.UniversalClient, logger logger.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err := client.Close(); err != nil {
				logger.Errorf("error in closing redis: %v", err)
			} else {
				logger.Info("redis closed gracefully")
			}

			return nil
		},
	})
}
