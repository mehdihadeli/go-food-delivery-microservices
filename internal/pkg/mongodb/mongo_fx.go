package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/health"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"mongofx",
	fx.Provide(
		provideConfig,
		NewMongoDB,
		fx.Annotate(
			NewMongoHealthChecker,
			fx.As(new(health.Health)),
			fx.ResultTags(fmt.Sprintf(`group:"%s"`, "healths")),
		),
	),
	fx.Invoke(registerHooks),
)

func registerHooks(lc fx.Lifecycle, client *mongo.Client, logger logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := client.Ping(ctx, nil)
			if err != nil {
				logger.Error("failed to ping mongo", zap.Error(err))
				return err
			}

			logger.Info("successfully pinged mongo")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := client.Disconnect(ctx); err != nil {
				logger.Errorf("error in disconnecting mongo: %v", err)
			} else {
				logger.Info("mongo disconnected gracefully")
			}

			return nil
		},
	})
}
