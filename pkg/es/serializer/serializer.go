package esSerializer

import (
	"encoding/json"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func SerializeToESEvent(aggregate contracts.IEventSourcedAggregateRoot, event interface{}, metadata *domain.Metadata) (*es.ESEvent, error) {
	eventBytes, err := serializer.Marshal(event)
	if err != nil {
		return nil, errors.Wrapf(err, "serializer.Marshal aggregateID: %s", aggregate.GetID())
	}

	metadataBytes, err := serializer.Marshal(metadata)
	if err != nil {
		return &es.ESEvent{}, errors.Wrapf(err, "serializer.Marshal aggregateID: %s", aggregate.GetID())
	}

	return es.NewESEventWIthMetadata(aggregate, es.EventType(typeMapper.GetTypeName(event)), eventBytes, metadataBytes), nil
}

func DeserializeGenericEventFromESEvent[T interface{}](event *es.ESEvent) (*T, error) {
	var targetEvent T
	if err := event.GetJsonData(&targetEvent); err != nil {
		return nil, errors.Wrapf(err, "event.GetJsonData type: %s", event.GetEventType())
	}
	return &targetEvent, nil
}

func DeserializeEventFromESEvent(event *es.ESEvent) (interface{}, error) {
	targetEventPointer := typeMapper.TypePointerInstanceByName(string(event.EventType))

	if err := event.GetJsonData(targetEventPointer); err != nil {
		return nil, errors.Wrapf(err, "event.GetJsonData type: %s", event.GetEventType())
	}

	return targetEventPointer, nil
}

func DeserializeEventFromRecordedEvent(recordedEvent *esdb.RecordedEvent) (interface{}, error) {
	targetEventPointer := typeMapper.TypePointerInstanceByName(recordedEvent.EventType)

	if err := json.Unmarshal(recordedEvent.Data, targetEventPointer); err != nil {
		return nil, errors.Wrapf(err, "event.GetJsonData type: %s", recordedEvent.EventType)
	}

	return targetEventPointer, nil
}

func DeserializeGenericEventFromRecordedEvent[T interface{}](recordedEvent *esdb.RecordedEvent) (T, error) {
	var targetEvent T

	if err := json.Unmarshal(recordedEvent.Data, &targetEvent); err != nil {
		return *new(T), errors.Wrapf(err, "event.GetJsonData type: %s", recordedEvent.EventType)
	}

	return targetEvent, nil
}

func DeserializeMetadataFromESEvent(event *es.ESEvent) (*domain.Metadata, error) {
	var targetMetadata domain.Metadata
	if err := event.GetJsonMetadata(&targetMetadata); err != nil {
		return nil, errors.Wrapf(err, "event.GetJsonMetadata")
	}
	return &targetMetadata, nil
}

func DeserializeMetadataFromRecordedEvent(recordedEvent *esdb.RecordedEvent) (*domain.Metadata, error) {
	var targetMetadata domain.Metadata
	if err := json.Unmarshal(recordedEvent.UserMetadata, &targetMetadata); err != nil {
		return nil, errors.Wrapf(err, "event.GetJsonMetadata")
	}
	return &targetMetadata, nil
}

func ToESEventFromRecorded(event *esdb.RecordedEvent) *es.ESEvent {
	id, err := uuid.FromString(event.EventID.String())
	if err != nil {
		return nil
	}

	return &es.ESEvent{
		EventID:   id,
		EventType: es.EventType(event.EventType),
		Data:      event.Data,
		Timestamp: event.CreatedDate,
		StreamID:  event.StreamID,
		Version:   int64(event.EventNumber),
		Metadata:  event.UserMetadata,
	}
}

func ToESEventFromEventData(event esdb.EventData) *es.ESEvent {
	id, err := uuid.FromString(event.EventID.String())
	if err != nil {
		return nil
	}

	return &es.ESEvent{
		EventID:   id,
		EventType: es.EventType(event.EventType),
		Data:      event.Data,
		Metadata:  event.Metadata,
	}
}

func ToESEventFromRecordedEvent(recordedEvent *esdb.RecordedEvent) (*es.ESEvent, error) {
	id, err := uuid.FromString(recordedEvent.EventID.String())
	if err != nil {
		return nil, err
	}

	return &es.ESEvent{
		EventID:   id,
		EventType: es.EventType(recordedEvent.EventType),
		Data:      recordedEvent.Data,
		Timestamp: recordedEvent.CreatedDate,
		StreamID:  recordedEvent.StreamID,
		Version:   int64(recordedEvent.Position.Commit),
		Metadata:  nil,
	}, nil
}

func ToEventData(e *es.ESEvent) esdb.EventData {
	return esdb.EventData{
		EventType:   string(e.EventType),
		ContentType: esdb.JsonContentType,
		Data:        e.Data,
		Metadata:    e.Metadata,
	}
}

func ToEventEnvelopeFromESEvent(event *es.ESEvent) (*domain.EventEnvelope, error) {
	deserializeEvent, err := DeserializeEventFromESEvent(event)
	if err != nil {
		return nil, err
	}

	deserializedMeta, err := DeserializeMetadataFromESEvent(event)
	if err != nil {
		return nil, err
	}

	return &domain.EventEnvelope{
		EventData: deserializeEvent,
		Metadata:  deserializedMeta,
	}, nil
}

func ToEventEnvelopeFromRecordedEvent(recordedEvent *esdb.RecordedEvent) (*domain.EventEnvelope, error) {

	deserializeEvent, err := DeserializeEventFromRecordedEvent(recordedEvent)
	if err != nil {
		return nil, err
	}

	deserializedMeta, err := DeserializeMetadataFromRecordedEvent(recordedEvent)
	if err != nil {
		return nil, err
	}

	return &domain.EventEnvelope{
		EventData: deserializeEvent,
		Metadata:  deserializedMeta,
	}, nil
}
