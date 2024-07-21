package health

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/health/contracts"
)

type healthService struct {
	healthParams contracts.HealthParams
}

func NewHealthService(
	healthParams contracts.HealthParams,
) contracts.HealthService {
	return &healthService{
		healthParams: healthParams,
	}
}

func (service *healthService) CheckHealth(ctx context.Context) contracts.Check {
	checks := make(contracts.Check)

	for _, health := range service.healthParams.Healths {
		checks[health.GetHealthName()] = contracts.NewStatus(
			health.CheckHealth(ctx),
		)
	}

	return checks
}
