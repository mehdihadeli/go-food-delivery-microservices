package eventstroredb

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gofrs/uuid"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es"
	appendResult "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/append_result"
	esSerializer "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/serializer"
	readPosition "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/stream_position/read_position"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/stream_position/truncatePosition"
	expectedStreamVersion "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/stream_version"
	esErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/eventstroredb/errors"
	"github.com/pkg/errors"
	uuid2 "github.com/satori/go.uuid"
	"io"
	"strings"
)

type EsdbSerializer struct {
	metadataSerializer esSerializer.MetadataSerializer
	eventSerializer    esSerializer.EventSerializer
}

func NewEsdbSerializer(metadataSerializer esSerializer.MetadataSerializer, eventSerializer esSerializer.EventSerializer) *EsdbSerializer {
	return &EsdbSerializer{
		metadataSerializer: metadataSerializer,
		eventSerializer:    eventSerializer,
	}
}

func (e *EsdbSerializer) StreamEventToEventData(streamEvent *es.StreamEvent) (esdb.EventData, error) {
	eventSerializationResult, err := e.eventSerializer.Serialize(streamEvent.Event)
	if err != nil {
		return *new(esdb.EventData), err
	}

	metadataSerializationResult, err := e.metadataSerializer.Serialize(streamEvent.Metadata)
	if err != nil {
		return *new(esdb.EventData), err
	}

	var contentType esdb.ContentType

	switch eventSerializationResult.ContentType {
	case "application/json":
		contentType = esdb.JsonContentType
	default:
		contentType = esdb.BinaryContentType
	}

	id, err := uuid.FromString(streamEvent.EventID.String())
	if err != nil {
		return *new(esdb.EventData), err
	}
	return esdb.EventData{
		EventID:     id,
		EventType:   eventSerializationResult.EventType,
		Data:        eventSerializationResult.Data,
		Metadata:    metadataSerializationResult,
		ContentType: contentType,
	}, nil
}

func (e *EsdbSerializer) ExpectedStreamVersionToEsdbExpectedRevision(expectedVersion expectedStreamVersion.ExpectedStreamVersion) esdb.ExpectedRevision {
	if expectedVersion.IsNoStream() {
		return esdb.NoStream{}
	}
	if expectedVersion.IsAny() {
		return esdb.Any{}
	}
	if expectedVersion.IsStreamExists() {
		return esdb.StreamExists{}
	}

	return esdb.StreamRevision{Value: uint64(expectedVersion.Value())}
}

func (e *EsdbSerializer) StreamReadPositionToStreamPosition(readPosition readPosition.StreamReadPosition) esdb.StreamPosition {
	if readPosition.IsEnd() {
		return esdb.End{}
	}
	if readPosition.IsStart() {
		return esdb.Start{}
	}

	return esdb.Revision(1)
}

func (e *EsdbSerializer) StreamTruncatePositionToInt64(truncatePosition truncatePosition.StreamTruncatePosition) uint64 {
	return uint64(truncatePosition.Value())
}

func (e *EsdbSerializer) EsdbReadStreamToResolvedEvents(stream *esdb.ReadStream) ([]*esdb.ResolvedEvent, error) {
	var events []*esdb.ResolvedEvent

	for {
		event, err := stream.Recv()
		if errors.Is(err, esdb.ErrStreamNotFound) {
			return nil, esErrors.NewStreamNotFoundError(err, event.Event.StreamID)
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, esErrors.NewReadStreamError(err)
		}

		events = append(events, event)
	}

	return events, nil
}

func (e *EsdbSerializer) EsdbPositionToStreamReadPosition(position esdb.Position) readPosition.StreamReadPosition {
	return readPosition.FromInt64(int64(position.Commit))
}

func (e *EsdbSerializer) ResolvedEventToStreamEvent(resolveEvent *esdb.ResolvedEvent) (*es.StreamEvent, error) {
	deserializedEvent, err := e.eventSerializer.Deserialize(resolveEvent.Event.Data, resolveEvent.Event.EventType, resolveEvent.Event.ContentType)
	if err != nil {
		return nil, err
	}

	deserializedMeta, err := e.metadataSerializer.Deserialize(resolveEvent.Event.UserMetadata)
	if err != nil {
		return nil, err
	}

	id, err := uuid2.FromString(resolveEvent.Event.EventID.String())
	if err != nil {
		return nil, err
	}

	return &es.StreamEvent{
		EventID:  id,
		Event:    deserializedEvent,
		Metadata: deserializedMeta,
		Version:  int64(resolveEvent.Event.EventNumber),
		Position: e.EsdbPositionToStreamReadPosition(resolveEvent.OriginalEvent().Position).Value(),
	}, nil
}

func (e *EsdbSerializer) ResolvedEventsToStreamEvents(resolveEvents []*esdb.ResolvedEvent) ([]*es.StreamEvent, error) {
	var streamEvents []*es.StreamEvent

	linq.From(resolveEvents).WhereT(func(item *esdb.ResolvedEvent) bool {
		return strings.HasPrefix(item.Event.EventType, "$") == false
	}).SelectT(func(item *esdb.ResolvedEvent) *es.StreamEvent {
		event, err := e.ResolvedEventToStreamEvent(item)
		if err != nil {
			return nil
		}
		return event
	}).ToSlice(&streamEvents)

	return streamEvents, nil
}

func (e *EsdbSerializer) EsdbWriteResultToAppendEventResult(writeResult *esdb.WriteResult) *appendResult.AppendEventsResult {
	return appendResult.From(writeResult.CommitPosition, writeResult.NextExpectedVersion)
}

func (e *EsdbSerializer) Serialize(data interface{}, metadata *core.Metadata) (*esdb.EventData, error) {
	serializedData, err := e.eventSerializer.SerializeObject(data)
	if err != nil {
		return nil, err
	}

	serializedMeta, err := e.metadataSerializer.Serialize(metadata)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewV4()
	return &esdb.EventData{
		EventID:     id,
		EventType:   serializedData.EventType,
		Data:        serializedData.Data,
		ContentType: esdb.JsonContentType,
		Metadata:    serializedMeta,
	}, nil
}

func (e *EsdbSerializer) Deserialize(resolveEvent *esdb.ResolvedEvent) (interface{}, *core.Metadata, error) {
	eventType := resolveEvent.Event.EventType
	data := resolveEvent.Event.Data
	meta := resolveEvent.Event.UserMetadata

	payload, err := e.eventSerializer.DeserializeObject(data, eventType, resolveEvent.Event.ContentType)
	if err != nil {
		return nil, nil, err
	}

	metadata, err := e.metadataSerializer.Deserialize(meta)
	if err != nil {
		return nil, nil, err
	}

	return payload, metadata, nil
}

func (e *EsdbSerializer) DomainEventToStreamEvent(domainEvent domain.IDomainEvent, meta *core.Metadata, position int64) *es.StreamEvent {
	return &es.StreamEvent{
		EventID:  uuid2.NewV4(),
		Event:    domainEvent,
		Metadata: meta,
		Version:  position,
		Position: position,
	}
}
