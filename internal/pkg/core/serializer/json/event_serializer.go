package json

import (
	"reflect"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/domain"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer/contratcs"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/typemapper"

	"emperror.dev/errors"
)

type DefaultEventJsonSerializer struct {
	serializer contratcs.Serializer
}

func NewDefaultEventJsonSerializer(serializer contratcs.Serializer) contratcs.EventSerializer {
	return &DefaultEventJsonSerializer{serializer: serializer}
}

func (s *DefaultEventJsonSerializer) Serialize(event domain.IDomainEvent) (*contratcs.EventSerializationResult, error) {
	return s.SerializeObject(event)
}

func (s *DefaultEventJsonSerializer) SerializeObject(event interface{}) (*contratcs.EventSerializationResult, error) {
	if event == nil {
		return &contratcs.EventSerializationResult{Data: nil, ContentType: s.ContentType()}, nil
	}

	// we use event short type name instead of full type name because this event in other receiver packages could have different package name
	eventType := typeMapper.GetTypeName(event)

	data, err := s.serializer.Marshal(event)
	if err != nil {
		return nil, errors.WrapIff(err, "error in Marshaling: `%s`", eventType)
	}

	result := &contratcs.EventSerializationResult{Data: data, ContentType: s.ContentType()}

	return result, nil
}

func (s *DefaultEventJsonSerializer) Deserialize(
	data []byte,
	eventType string,
	contentType string,
) (domain.IDomainEvent, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.EmptyInstanceByTypeNameAndImplementedInterface[domain.IDomainEvent](
		eventType,
	)

	if targetEventPointer == nil {
		return nil, errors.Errorf("event type `%s` is not impelemted IDomainEvent or can't be instansiated", eventType)
	}

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := s.serializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.WrapIff(err, "error in Unmarshaling: `%s`", eventType)
	}

	return targetEventPointer.(domain.IDomainEvent), nil
}

func (s *DefaultEventJsonSerializer) DeserializeObject(
	data []byte,
	eventType string,
	contentType string,
) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.InstanceByTypeName(eventType)

	if targetEventPointer == nil {
		return nil, errors.Errorf("event type `%s` can't be instansiated", eventType)
	}

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := s.serializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.WrapIff(err, "error in Unmarshaling: `%s`", eventType)
	}

	return targetEventPointer, nil
}

func (s *DefaultEventJsonSerializer) DeserializeType(
	data []byte,
	eventType reflect.Type,
	contentType string,
) (domain.IDomainEvent, error) {
	if data == nil {
		return nil, nil
	}

	// we use event short type name instead of full type name because this event in other receiver packages could have different package name
	eventTypeName := typeMapper.GetTypeName(eventType)

	return s.Deserialize(data, eventTypeName, contentType)
}

func (s *DefaultEventJsonSerializer) ContentType() string {
	return "application/json"
}

func (s *DefaultEventJsonSerializer) Serializer() contratcs.Serializer {
	return s.serializer
}
