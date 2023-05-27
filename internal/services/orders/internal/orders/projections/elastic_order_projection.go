package projections

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/es/contracts/projection"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/es/models"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/contracts/repositories"
)

type elasticOrderProjection struct {
	elasticOrderReadRepository repositories.OrderReadRepository
}

func NewElasticOrderProjection(elasticOrderReadRepository repositories.OrderReadRepository) projection.IProjection {
	return &elasticOrderProjection{elasticOrderReadRepository: elasticOrderReadRepository}
}

func (e elasticOrderProjection) ProcessEvent(ctx context.Context, streamEvent *models.StreamEvent) error {
	// TODO: Handling and projecting event to elastic read model
	return nil
}
