package json

import (
	"reflect"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer/contratcs"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/typemapper"

	"emperror.dev/errors"
)

type DefaultMessageJsonSerializer struct {
	serializer contratcs.Serializer
}

func NewDefaultMessageJsonSerializer(serializer contratcs.Serializer) contratcs.MessageSerializer {
	return &DefaultMessageJsonSerializer{serializer: serializer}
}

func (m *DefaultMessageJsonSerializer) Serialize(message types.IMessage) (*contratcs.EventSerializationResult, error) {
	return m.SerializeObject(message)
}

func (m *DefaultMessageJsonSerializer) SerializeObject(
	message interface{},
) (*contratcs.EventSerializationResult, error) {
	if message == nil {
		return &contratcs.EventSerializationResult{Data: nil, ContentType: m.ContentType()}, nil
	}

	// we use message short type name instead of full type name because this message in other receiver packages could have different package name
	eventType := typeMapper.GetTypeName(message)

	data, err := m.serializer.Marshal(message)
	if err != nil {
		return nil, errors.WrapIff(err, "error in Marshaling: `%s`", eventType)
	}

	result := &contratcs.EventSerializationResult{Data: data, ContentType: m.ContentType()}

	return result, nil
}

func (m *DefaultMessageJsonSerializer) SerializeEnvelop(
	messageEnvelop types.MessageEnvelope,
) (*contratcs.EventSerializationResult, error) {
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

func (m *DefaultMessageJsonSerializer) Serializer() contratcs.Serializer {
	return m.serializer
}
