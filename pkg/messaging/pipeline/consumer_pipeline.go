package pipeline

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
)

// ConsumerHandlerFunc is a continuation for the next task to execute in the pipeline
type ConsumerHandlerFunc func() error

// ConsumerPipeline is a Pipeline for wrapping the inner consumer handler.
type ConsumerPipeline interface {
	Handle(ctx context.Context, consumerContext types.MessageConsumeContext, next ConsumerHandlerFunc) error
}
