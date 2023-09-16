package redis

import (
	"context"
	"fmt"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/health"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("redisfx",
	fx.Provide(
		NewRedisClient,
		func(client *redis.Client) redis.UniversalClient {
			return client
		},
		// will create new instance of redis client
		//fx.Annotate(
		//	NewRedisClient,
		//	fx.As(new(redis.UniversalClient)),
		//),
		fx.Annotate(
			NewRedisHealthChecker,
			fx.As(new(health.Health)),
			fx.ResultTags(fmt.Sprintf(`group:"%s"`, "healths")),
		),
		provideConfig),
	fx.Invoke(registerHooks),
)

func registerHooks(lc fx.Lifecycle, client redis.UniversalClient, logger logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Ping(ctx).Err()
		},
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
