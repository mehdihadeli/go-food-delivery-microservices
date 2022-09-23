package json

import (
	"emperror.dev/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer/jsonSerializer"
	"reflect"
)

type JsonEventSerializer struct {
}

func NewJsonEventSerializer() *JsonEventSerializer {
	return &JsonEventSerializer{}
}

func (s *JsonEventSerializer) Serialize(event interface{}) (*serializer.EventSerializationResult, error) {
	if event == nil {
		return &serializer.EventSerializationResult{Data: nil, ContentType: s.ContentType()}, nil
	}

	// here we just get type name instead of full type name
	eventType := typeMapper.GetTypeName(event)

	data, err := jsonSerializer.Marshal(event)
	if err != nil {
		return nil, errors.WrapIff(err, "event.GetJsonData type: %s", eventType)
	}

	result := &serializer.EventSerializationResult{Data: data, ContentType: s.ContentType()}

	return result, nil
}

func (s *JsonEventSerializer) Deserialize(data []byte, eventType string, contentType string) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.InstanceByTypeName(eventType)

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := jsonSerializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.WrapIff(err, "event.GetJsonData type: %s", eventType)
	}

	return targetEventPointer, nil
}

func (s *JsonEventSerializer) DeserializeType(data []byte, eventType reflect.Type, contentType string) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.InstanceByType(eventType)

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := jsonSerializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.WrapIff(err, "event.GetJsonData type: %s", eventType)
	}

	return targetEventPointer, nil
}

func (s *JsonEventSerializer) DeserializeMessage(data []byte, eventType string, contentType string) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.InstanceByTypeNameAndImplementedInterface[types.IMessage](eventType)

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := jsonSerializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.WrapIff(err, "event.GetJsonData type: %s", eventType)
	}

	return targetEventPointer, nil
}

func (s *JsonEventSerializer) DeserializeEvent(data []byte, eventType string, contentType string) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.InstanceByTypeNameAndImplementedInterface[core.IEvent](eventType)

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := jsonSerializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.WrapIff(err, "event.GetJsonData type: %s", eventType)
	}

	return targetEventPointer, nil
}

func (s *JsonEventSerializer) ContentType() string {
	return "application/json"
}
