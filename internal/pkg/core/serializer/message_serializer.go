package serializer

import (
	"reflect"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/types"
)

type MessageSerializer interface {
	Serialize(message types.IMessage) (*EventSerializationResult, error)
	SerializeObject(message interface{}) (*EventSerializationResult, error)
	SerializeEnvelop(messageEnvelop types.MessageEnvelope) (*EventSerializationResult, error)
	Deserialize(data []byte, messageType string, contentType string) (types.IMessage, error)
	DeserializeObject(data []byte, messageType string, contentType string) (interface{}, error)
	DeserializeType(data []byte, messageType reflect.Type, contentType string) (types.IMessage, error)
	ContentType() string
	Serializer() Serializer
}
