package contracts

import (
	"context"
	"testing"

	postgres "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgrespgx"
)

type PostgresPgxContainer interface {
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*PostgresContainerOptions,
	) (*postgres.PostgresPgxOptions, error)
	Cleanup(ctx context.Context) error
}
