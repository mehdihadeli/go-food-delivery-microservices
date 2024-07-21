package contracts

import (
	"context"
	"testing"

	redis2 "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/redis"
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
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*RedisContainerOptions,
	) (*redis2.RedisOptions, error)
	Cleanup(ctx context.Context) error
}
