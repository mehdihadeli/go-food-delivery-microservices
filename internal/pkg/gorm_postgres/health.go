package gormPostgres

import (
	"context"
	"database/sql"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/health"
)

type gormHealthChecker struct {
	client *sql.DB
}

func NewGormHealthChecker(client *sql.DB) health.Health {
	return &gormHealthChecker{client}
}

func (healthChecker *gormHealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.PingContext(ctx)
}

func (healthChecker *gormHealthChecker) GetHealthName() string {
	return "postgres"
}
