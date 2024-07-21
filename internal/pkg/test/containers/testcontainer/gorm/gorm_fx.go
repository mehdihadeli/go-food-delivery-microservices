package gorm

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	gormPostgres "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm"
)

var GormContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *gormPostgres.GormOptions, logger logger.Logger) (*gormPostgres.GormOptions, error) {
		return NewGormTestContainers(logger).PopulateContainerOptions(ctx, t)
	}
}
