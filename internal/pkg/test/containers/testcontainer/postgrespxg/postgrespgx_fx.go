package postgrespxg

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	postgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgres_pgx"
)

var PostgresPgxContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *postgres.PostgresPgxOptions, logger logger.Logger) (*postgres.PostgresPgxOptions, error) {
		return NewPostgresPgxContainers(logger).CreatingContainerOptions(ctx, t)
	}
}
