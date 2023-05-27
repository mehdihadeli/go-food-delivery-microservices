package domain

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/metadata"
)

type EventEnvelope struct {
	EventData interface{}
	Metadata  metadata.Metadata
}
