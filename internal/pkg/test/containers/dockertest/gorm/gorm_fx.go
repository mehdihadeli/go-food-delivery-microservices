package gorm

import (
	"context"
	"testing"

	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/postgresGorm"
)

var GormDockerTestConatnerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *gormPostgres.GormOptions) (*gormPostgres.GormOptions, error) {
		return NewGormDockerTest().PopulateContainerOptions(ctx, t)
	}
}
