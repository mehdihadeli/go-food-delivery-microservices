package domain

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/metadata"
)

type EventEnvelope struct {
	EventData interface{}
	Metadata  metadata.Metadata
}
