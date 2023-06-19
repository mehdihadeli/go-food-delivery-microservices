package serializer

import (
	"reflect"

	"emperror.dev/errors"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/events"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"
)

type EventSerializer interface {
	Serialize(event interface{}) (*EventSerializationResult, error)
	Deserialize(data []byte, eventType string, contentType string) (interface{}, error)
	DeserializeType(data []byte, eventType reflect.Type, contentType string) (interface{}, error)
	DeserializeMessage(data []byte, eventType string, contentType string) (interface{}, error)
	DeserializeEvent(data []byte, eventType string, contentType string) (interface{}, error)
	ContentType() string
	Serializer() Serializer
}

type EventSerializationResult struct {
	Data        []byte
	ContentType string
}

type DefaultEventSerializer struct {
	serializer Serializer
}

func NewDefaultEventSerializer(serializer Serializer) EventSerializer {
	return &DefaultEventSerializer{serializer: serializer}
}

func (s *DefaultEventSerializer) Serializer() Serializer {
	return s.serializer
}

func (s *DefaultEventSerializer) Serialize(
	event interface{},
) (*EventSerializationResult, error) {
	if event == nil {
		return &EventSerializationResult{Data: nil, ContentType: s.ContentType()}, nil
	}

	// here we just get type name instead of full type name
	eventType := typeMapper.GetTypeName(event)

	data, err := s.serializer.Marshal(event)
	if err != nil {
		return nil, errors.WrapIff(err, "event.GetJsonData type: %s", eventType)
	}

	result := &EventSerializationResult{Data: data, ContentType: s.ContentType()}

	return result, nil
}

func (s *DefaultEventSerializer) Deserialize(
	data []byte,
	eventType string,
	contentType string,
) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.InstanceByTypeName(eventType)

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := s.serializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.WrapIff(err, "event.GetJsonData type: %s", eventType)
	}

	return targetEventPointer, nil
}

func (s *DefaultEventSerializer) DeserializeType(
	data []byte,
	eventType reflect.Type,
	contentType string,
) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.InstanceByType(eventType)

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := s.serializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.WrapIff(err, "event.GetJsonData type: %s", eventType)
	}

	return targetEventPointer, nil
}

func (s *DefaultEventSerializer) DeserializeMessage(
	data []byte,
	eventType string,
	contentType string,
) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.InstanceByTypeNameAndImplementedInterface[types.IMessage](
		eventType,
	)

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := s.serializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.WrapIff(err, "event.GetJsonData type: %s", eventType)
	}

	return targetEventPointer, nil
}

func (s *DefaultEventSerializer) DeserializeEvent(
	data []byte,
	eventType string,
	contentType string,
) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.InstanceByTypeNameAndImplementedInterface[events.IEvent](
		eventType,
	)

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := s.serializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.WrapIff(err, "event.GetJsonData type: %s", eventType)
	}

	return targetEventPointer, nil
}

func (s *DefaultEventSerializer) ContentType() string {
	return "application/json"
}
