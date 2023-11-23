package contracts

import (
	"context"
	"testing"

	postgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgres_pgx"
)

type PostgresPgxContainer interface {
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*PostgresContainerOptions,
	) (*postgres.PostgresPgxOptions, error)
	Cleanup(ctx context.Context) error
}
