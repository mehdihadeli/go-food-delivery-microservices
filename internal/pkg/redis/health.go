package redis

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/health/contracts"

	"github.com/redis/go-redis/v9"
)

type RedisHealthChecker struct {
	client *redis.Client
}

func NewRedisHealthChecker(client *redis.Client) contracts.Health {
	return &RedisHealthChecker{client}
}

func (healthChecker *RedisHealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.Ping(ctx).Err()
}

func (healthChecker *RedisHealthChecker) GetHealthName() string {
	return "redis"
}
