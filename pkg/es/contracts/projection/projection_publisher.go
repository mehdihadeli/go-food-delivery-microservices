package projection

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models"
)

type IProjectionPublisher interface {
	Publish(ctx context.Context, streamEvent *models.StreamEvent) error
}
