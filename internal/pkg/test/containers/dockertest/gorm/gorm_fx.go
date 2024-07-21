package gorm

import (
	"context"
	"testing"

	gormPostgres "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm"
)

var GormDockerTestConatnerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *gormPostgres.GormOptions) (*gormPostgres.GormOptions, error) {
		return NewGormDockerTest().PopulateContainerOptions(ctx, t)
	}
}
