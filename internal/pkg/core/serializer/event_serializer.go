package serializer

import (
	"reflect"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/domain"
)

type EventSerializer interface {
	Serialize(event domain.IDomainEvent) (*EventSerializationResult, error)
	SerializeObject(event interface{}) (*EventSerializationResult, error)
	Deserialize(data []byte, eventType string, contentType string) (domain.IDomainEvent, error)
	DeserializeObject(data []byte, eventType string, contentType string) (interface{}, error)
	DeserializeType(data []byte, eventType reflect.Type, contentType string) (domain.IDomainEvent, error)
	ContentType() string
	Serializer() Serializer
}
