package contracts

import "context"

type Health interface {
	CheckHealth(ctx context.Context) error
	GetHealthName() string
}

type HealthService interface {
	CheckHealth(ctx context.Context) Check
}
