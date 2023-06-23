package gorm

import (
	"context"
	"testing"

	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
)

var GormContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *gormPostgres.GormOptions) (*gormPostgres.GormOptions, error) {
		return NewGormTestContainers().CreatingContainerOptions(ctx, t)
	}
}
