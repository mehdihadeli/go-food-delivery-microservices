package gorm

import (
	"context"
	"testing"

	gormPostgres "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/gorm_postgres"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
)

var GormContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *gormPostgres.GormOptions, logger logger.Logger) (*gormPostgres.GormOptions, error) {
		return NewGormTestContainers(logger).CreatingContainerOptions(ctx, t)
	}
}
