package mongo

import (
	"context"
	"testing"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mongodb"
)

var MongoContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *mongodb.MongoDbOptions) (*mongodb.MongoDbOptions, error) {
		return NewMongoTestContainers().CreatingContainerOptions(ctx, t)
	}
}
