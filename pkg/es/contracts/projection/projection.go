package projection

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models"
)

type IProjection interface {
	ProcessEvent(ctx context.Context, streamEvent *models.StreamEvent) error
}
