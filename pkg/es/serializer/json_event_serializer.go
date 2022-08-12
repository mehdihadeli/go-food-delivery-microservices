package esSerializer

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	esSerializer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/serializer"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer/jsonSerializer"
	"github.com/pkg/errors"
)

type JsonEventSerializer struct {
}

func NewJsonEventSerializer() *JsonEventSerializer {
	return &JsonEventSerializer{}
}

func (s *JsonEventSerializer) Serialize(event domain.IDomainEvent) (*esSerializer.EventSerializationResult, error) {
	eventType := typeMapper.GetTypeName(event)

	data, err := jsonSerializer.Marshal(event)
	if err != nil {
		return nil, errors.Wrapf(err, "event.GetJsonData type: %s", eventType)
	}

	result := &esSerializer.EventSerializationResult{Data: data, ContentType: s.ContentType(), EventType: eventType}

	return result, nil
}

func (s *JsonEventSerializer) Deserialize(data []byte, eventType string, contentType string) (domain.IDomainEvent, error) {
	targetEventPointer := typeMapper.InstancePointerByTypeName(eventType)

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := jsonSerializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.Wrapf(err, "event.GetJsonData type: %s", eventType)
	}

	return targetEventPointer.(domain.IDomainEvent), nil
}

func (s *JsonEventSerializer) ContentType() string {
	return "application/json"
}
