package json

import (
	"emperror.dev/errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/serializer"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer/jsonSerializer"
)

type JsonEventSerializer struct {
}

func NewJsonEventSerializer() *JsonEventSerializer {
	return &JsonEventSerializer{}
}

func (s *JsonEventSerializer) Serialize(event interface{}) (*serializer.EventSerializationResult, error) {
	if event == nil {
		return &serializer.EventSerializationResult{Data: nil, ContentType: s.ContentType(), EventType: ""}, nil
	}

	eventType := typeMapper.GetTypeName(event)

	data, err := jsonSerializer.Marshal(event)
	if err != nil {
		return nil, errors.WrapIff(err, "event.GetJsonData type: %s", eventType)
	}

	result := &serializer.EventSerializationResult{Data: data, ContentType: s.ContentType(), EventType: eventType}

	return result, nil
}

func (s *JsonEventSerializer) Deserialize(data []byte, eventType string, contentType string) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.InstancePointerByTypeName(eventType)

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
