package projection

import (
    "context"

    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/es/models"
)

type IProjection interface {
    ProcessEvent(ctx context.Context, streamEvent *models.StreamEvent) error
}
