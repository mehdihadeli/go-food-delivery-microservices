package health

import (
	"context"
)

type HealthService interface {
	CheckHealth(ctx context.Context) Check
}

type healthService struct {
	healthParams HealthParams
}

func NewHealthService(
	healthParams HealthParams,
) HealthService {
	return &healthService{
		healthParams: healthParams,
	}
}

func (service *healthService) CheckHealth(ctx context.Context) Check {
	checks := make(Check)

	for _, health := range service.healthParams.Healths {
		checks[health.GetHealthName()] = NewStatus(health.CheckHealth(ctx))
	}

	return checks
}
