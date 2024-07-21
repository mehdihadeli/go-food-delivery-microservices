package json

import (
	"reflect"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/messaging/types"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/serializer"
	typeMapper "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/reflection/typemapper"

	"emperror.dev/errors"
)

type DefaultMessageJsonSerializer struct {
	serializer serializer.Serializer
}

func NewDefaultMessageJsonSerializer(s serializer.Serializer) serializer.MessageSerializer {
	return &DefaultMessageJsonSerializer{serializer: s}
}

func (m *DefaultMessageJsonSerializer) Serialize(message types.IMessage) (*serializer.EventSerializationResult, error) {
	return m.SerializeObject(message)
}

func (m *DefaultMessageJsonSerializer) SerializeObject(
	message interface{},
) (*serializer.EventSerializationResult, error) {
	if message == nil {
		return &serializer.EventSerializationResult{Data: nil, ContentType: m.ContentType()}, nil
	}

	// we use message short type name instead of full type name because this message in other receiver packages could have different package name
	eventType := typeMapper.GetTypeName(message)

	data, err := m.serializer.Marshal(message)
	if err != nil {
		return nil, errors.WrapIff(err, "error in Marshaling: `%s`", eventType)
	}

	result := &serializer.EventSerializationResult{Data: data, ContentType: m.ContentType()}

	return result, nil
}

func (m *DefaultMessageJsonSerializer) SerializeEnvelop(
	messageEnvelop types.MessageEnvelope,
) (*serializer.EventSerializationResult, error) {
	// TODO implement me
	panic("implement me")
}

func (m *DefaultMessageJsonSerializer) Deserialize(
	data []byte,
	messageType string,
	contentType string,
) (types.IMessage, error) {
	if data == nil {
		return nil, nil
	}

	targetMessagePointer := typeMapper.EmptyInstanceByTypeNameAndImplementedInterface[types.IMessage](
		messageType,
	)

	if targetMessagePointer == nil {
		return nil, errors.Errorf("message type `%s` is not impelemted IMessage or can't be instansiated", messageType)
	}

	if contentType != m.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := m.serializer.Unmarshal(data, targetMessagePointer); err != nil {
		return nil, errors.WrapIff(err, "error in Unmarshaling: `%s`", messageType)
	}

	return targetMessagePointer.(types.IMessage), nil
}

func (m *DefaultMessageJsonSerializer) DeserializeObject(
	data []byte,
	messageType string,
	contentType string,
) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetMessagePointer := typeMapper.InstanceByTypeName(messageType)

	if targetMessagePointer == nil {
		return nil, errors.Errorf("message type `%s` can't be instansiated", messageType)
	}

	if contentType != m.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := m.serializer.Unmarshal(data, targetMessagePointer); err != nil {
		return nil, errors.WrapIff(err, "error in Unmarshaling: `%s`", messageType)
	}

	return targetMessagePointer, nil
}

func (m *DefaultMessageJsonSerializer) DeserializeType(
	data []byte,
	messageType reflect.Type,
	contentType string,
) (types.IMessage, error) {
	if data == nil {
		return nil, nil
	}

	// we use message short type name instead of full type name because this message in other receiver packages could have different package name
	messageTypeName := typeMapper.GetTypeName(messageType)

	return m.Deserialize(data, messageTypeName, contentType)
}

func (m *DefaultMessageJsonSerializer) ContentType() string {
	return "application/json"
}

func (m *DefaultMessageJsonSerializer) Serializer() serializer.Serializer {
	return m.serializer
}
