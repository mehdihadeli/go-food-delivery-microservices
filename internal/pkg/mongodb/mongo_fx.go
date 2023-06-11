package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module(
	"mongofx",
	fx.Provide(
		provideConfig,
		NewMongoDB,
	),
	fx.Invoke(registerHooks),
)

func registerHooks(lc fx.Lifecycle, client *mongo.Client, logger logger.Logger) {
	lc.Append(fx.Hook{
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
