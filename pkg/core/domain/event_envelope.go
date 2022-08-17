package domain

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"

type EventEnvelope struct {
	EventData interface{}
	Metadata  *core.Metadata
}
