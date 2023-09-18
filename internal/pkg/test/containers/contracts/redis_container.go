package contracts

import (
	"context"
	"testing"

	redis2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/redis"

	"github.com/redis/go-redis/v9"
)

type RedisContainerOptions struct {
	Host      string
	Port      string
	HostPort  int
	Database  int
	ImageName string
	Name      string
	Tag       string
	PoolSize  int
}

type RedisContainer interface {
	Start(
		ctx context.Context,
		t *testing.T,
		options ...*RedisContainerOptions,
	) (redis.UniversalClient, error)
	CreatingContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*RedisContainerOptions,
	) (*redis2.RedisOptions, error)
	Cleanup(ctx context.Context) error
}
