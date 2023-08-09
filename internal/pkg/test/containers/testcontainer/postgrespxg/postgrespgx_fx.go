package postgrespxg

import (
	"context"
	"testing"

	postgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgres_pgx"
)

var PostgresPgxContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *postgres.PostgresPgxOptions) (*postgres.PostgresPgxOptions, error) {
		return NewPostgresPgxContainers().CreatingContainerOptions(ctx, t)
	}
}
