package domain

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/metadata"
)

type EventEnvelope struct {
	EventData interface{}
	Metadata  metadata.Metadata
}
