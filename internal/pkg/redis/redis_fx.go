package redis

import (
	"context"
	"fmt"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/health"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var (
	// Module provided to fxlog
	// https://uber-go.github.io/fx/modules.html
	Module = fx.Module("redisfx", redisProviders, redisInvokes) //nolint:gochecknoglobals

	redisProviders = fx.Options(fx.Provide( //nolint:gochecknoglobals
		NewRedisClient,
		func(client *redis.Client) redis.UniversalClient {
			return client
		},
		//// will create new instance of redis client instead of reusing current instance of `redis.Client`
		//fx.Annotate(
		//	NewRedisClient,
		//	fx.As(new(redis.UniversalClient)),
		//),
		fx.Annotate(
			NewRedisHealthChecker,
			fx.As(new(health.Health)),
			fx.ResultTags(fmt.Sprintf(`group:"%s"`, "healths")),
		),
		provideConfig))

	redisInvokes = fx.Options(fx.Invoke(registerHooks), //nolint:gochecknoglobals
		fx.Invoke(EnableTracing))
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

func EnableTracing(redis *redis.Client) error {
	return redisotel.InstrumentTracing(redis)
}
