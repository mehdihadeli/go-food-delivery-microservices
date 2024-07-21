package contracts

import (
	"context"

	"gorm.io/gorm"
)

type GormDBContext interface {
	WithTx(ctx context.Context) (GormDBContext, error)
	WithTxIfExists(ctx context.Context) GormDBContext
	RunInTx(ctx context.Context, action ActionFunc) error
	DB() *gorm.DB
}

type ActionFunc func(ctx context.Context, gormContext GormDBContext) error
