package es

import (
	"context"
	"emperror.dev/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/projection"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models"
)

type projectionPublisher struct {
	projections []projection.IProjection
}

func NewProjectionPublisher(projections []projection.IProjection) projection.IProjectionPublisher {
	return &projectionPublisher{projections: projections}
}

func (p projectionPublisher) Publish(ctx context.Context, streamEvent *models.StreamEvent) error {
	if streamEvent == nil {
		return nil
	}

	if p.projections == nil {
		return nil
	}

	for _, projection := range p.projections {
		err := projection.ProcessEvent(ctx, streamEvent)
		if err != nil {
			return errors.WrapIf(err, "error in processing projection")
		}
	}

	return nil
}
