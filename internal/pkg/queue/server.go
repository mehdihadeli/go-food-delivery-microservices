package queue

import (
	"context"
	"fmt"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	redis2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/redis"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

func NewServeMux() *asynq.ServeMux {
	return asynq.NewServeMux()
}

func NewServer(config *redis2.RedisOptions, logger logger.Logger) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%d", config.Host, config.Port)},
		asynq.Config{Concurrency: 10},
	)
}

func HookServer(lifecycle fx.Lifecycle, server *asynq.Server, mux *asynq.ServeMux) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.Run(mux); err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.Shutdown()
			return nil
		},
	})
}
