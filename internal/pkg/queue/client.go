package queue

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"

	redis2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/redis"
)

func NewClient(config *redis2.RedisOptions) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%d", config.Host, config.Port)})
}

func HookClient(lifecycle fx.Lifecycle, client *asynq.Client) {
	lifecycle.Append(fx.Hook{
		OnStop: func(c context.Context) error {
			return client.Close()
		},
	})
}
