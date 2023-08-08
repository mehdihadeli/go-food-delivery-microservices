package models

import (
	uuid "github.com/satori/go.uuid"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/domain"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/metadata"
)

type StreamEvent struct {
	EventID  uuid.UUID
	Version  int64
	Position int64
	Event    domain.IDomainEvent
	Metadata metadata.Metadata
}
