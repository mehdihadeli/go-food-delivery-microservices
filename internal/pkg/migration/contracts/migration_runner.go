package contracts

import "context"

type PostgresMigrationRunner interface {
	Up(ctx context.Context, version uint) error
	Down(ctx context.Context, version uint) error
}
