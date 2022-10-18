package data

import (
	"context"
)

type UnitOfWorkActionFunc func() error

type UnitOfWork interface {
	// SaveWithTx executes the given UnitOfWorkActionFunc atomically (inside a DB transaction).
	SaveWithTx(ctx context.Context, action UnitOfWorkActionFunc) error
}
