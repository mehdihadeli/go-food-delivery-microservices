package redis

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/redis"
)

var RedisContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *redis.RedisOptions) (*redis.RedisOptions, error) {
		return NewRedisTestContainers().CreatingContainerOptions(ctx, t)
	}
}
