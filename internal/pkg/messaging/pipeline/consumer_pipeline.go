package pipeline

import (
    "context"

    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
)

// ConsumerHandlerFunc is a continuation for the next task to execute in the pipeline
type ConsumerHandlerFunc func() error

// ConsumerPipeline is a Pipeline for wrapping the inner consumer handler.
type ConsumerPipeline interface {
    Handle(ctx context.Context, consumerContext types.MessageConsumeContext, next ConsumerHandlerFunc) error
}
