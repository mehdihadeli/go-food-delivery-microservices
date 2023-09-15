package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/health"
)

type mongoHealthChecker struct {
	client *mongo.Client
}

func NewMongoHealthChecker(client *mongo.Client) health.Health {
	return &mongoHealthChecker{client}
}

func (healthChecker *mongoHealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.Ping(ctx, nil)
}

func (healthChecker *mongoHealthChecker) GetHealthName() string {
	return "mongodb"
}
