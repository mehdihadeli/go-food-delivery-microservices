package contracts

import (
	"context"

	"gorm.io/gorm"
)

type IGormDBContext interface {
	WithTx(ctx context.Context) (IGormDBContext, error)
	WithTxIfExists(ctx context.Context) IGormDBContext
	RunInTx(ctx context.Context, action ActionFunc) error
	DB() *gorm.DB
}

type ActionFunc func(ctx context.Context, gormContext IGormDBContext) error
