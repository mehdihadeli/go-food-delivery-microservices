package consumer

import (
    "context"

    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
)

type ConsumerHandler interface {
    Handle(ctx context.Context, consumeContext types.MessageConsumeContext) error
}
