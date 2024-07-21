package projection

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/es/models"
)

type IProjectionPublisher interface {
	Publish(ctx context.Context, streamEvent *models.StreamEvent) error
}
