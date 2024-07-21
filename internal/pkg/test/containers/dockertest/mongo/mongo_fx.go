package mongo

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mongodb"
)

var MongoDockerTestContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *mongodb.MongoDbOptions, logger logger.Logger) (*mongodb.MongoDbOptions, error) {
		return NewMongoDockerTest().PopulateContainerOptions(ctx, t)
	}
}
