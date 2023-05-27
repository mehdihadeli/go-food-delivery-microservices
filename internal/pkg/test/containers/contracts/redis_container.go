package contracts

import (
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
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
	Start(ctx context.Context, t *testing.T, options ...*RedisContainerOptions) (redis.UniversalClient, error)
	Cleanup(ctx context.Context) error
}
