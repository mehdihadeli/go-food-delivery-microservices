package es

import (
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts"
	streamName "github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/stream_name"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/serializer"
	uuid "github.com/satori/go.uuid"
	"time"
)

// EventType is the type of any event, used as its unique identifier.
type EventType string

// ESEvent is an internal representation of an event, returned when the Aggregate
// uses NewEvent to create a new event. The events loaded from the db is
// represented by each DBs internal event type, implementing Event.
type ESEvent struct {
	EventID       uuid.UUID
	EventType     EventType
	Data          []byte
	Timestamp     time.Time
	AggregateType types.AggregateType
	StreamID      string
	Version       int64
	Metadata      []byte
}

// NewESEvent new base Event constructor with configured EventID, Aggregate properties and Timestamp.
func NewESEvent(aggregate contracts.IEventSourcedAggregateRoot, eventType EventType) *ESEvent {
	return &ESEvent{
		EventID:       uuid.NewV4(),
		AggregateType: aggregate.GetType(),
		StreamID:      streamName.For(aggregate),
		Version:       aggregate.GetVersion(),
		EventType:     eventType,
		Timestamp:     time.Now(),
	}
}

func NewESEventWIthMetadata(aggregate contracts.IEventSourcedAggregateRoot, eventType EventType, data []byte, metadata []byte) *ESEvent {
	return &ESEvent{
		EventID:       uuid.NewV4(),
		StreamID:      streamName.For(aggregate),
		EventType:     eventType,
		AggregateType: aggregate.GetType(),
		Version:       aggregate.GetVersion(),
		Data:          data,
		Metadata:      metadata,
		Timestamp:     time.Now(),
	}
}

// GetEventID get EventID of the Event.
func (e *ESEvent) GetEventID() uuid.UUID {
	return e.EventID
}

// GetTimeStamp get timestamp of the Event.
func (e *ESEvent) GetTimeStamp() time.Time {
	return e.Timestamp
}

// GetData The data attached to the Event serialized to bytes.
func (e *ESEvent) GetData() []byte {
	return e.Data
}

// SetData add the data attached to the Event serialized to bytes.
func (e *ESEvent) SetData(data []byte) *ESEvent {
	e.Data = data
	return e
}

// GetJsonData json unmarshal data attached to the Event.
func (e *ESEvent) GetJsonData(data interface{}) error {
	err := serializer.Unmarshal(e.GetData(), data)
	if err != nil {
		return err
	}

	return nil
}

// SetJsonData serialize to json and set data attached to the Event.
func (e *ESEvent) SetJsonData(data interface{}) error {
	dataBytes, err := serializer.Marshal(data)
	if err != nil {
		return err
	}

	e.Data = dataBytes
	return nil
}

// GetMetadata is app-specific metadata such as request ID, originating user etc.
func (e *ESEvent) GetMetadata() []byte {
	return e.Metadata
}

// SetMetadata add app-specific metadata serialized as json for the Event.
func (e *ESEvent) SetMetadata(metaData interface{}) error {

	metaDataBytes, err := serializer.Marshal(metaData)
	if err != nil {
		return err
	}

	e.Metadata = metaDataBytes
	return nil
}

// GetJsonMetadata unmarshal app-specific metadata serialized as json for the Event.
func (e *ESEvent) GetJsonMetadata(metaData interface{}) error {
	//https://www.sohamkamani.com/golang/json/#decoding-json-to-maps---unstructured-data
	err := serializer.Unmarshal(e.GetMetadata(), metaData)
	if err != nil {
		return err
	}

	return nil
}

// GetEventType returns the EventType of the event.
func (e *ESEvent) GetEventType() EventType {
	return e.EventType
}

// GetAggregateType is the AggregateType that the Event can be applied to.
func (e *ESEvent) GetAggregateType() types.AggregateType {
	return e.AggregateType
}

// SetAggregateType set the AggregateType that the Event can be applied to.
func (e *ESEvent) SetAggregateType(aggregateType types.AggregateType) {
	e.AggregateType = aggregateType
}

// GetVersion is the version of the Aggregate after the Event has been applied.
func (e *ESEvent) GetVersion() int64 {
	return e.Version
}

// SetVersion set the version of the Aggregate.
func (e *ESEvent) SetVersion(aggregateVersion int64) {
	e.Version = aggregateVersion
}

// GetString A string representation of the Event.
func (e *ESEvent) GetString() string {
	return fmt.Sprintf("event: %+v", e)
}

func (e *ESEvent) String() string {
	return fmt.Sprintf("(Event): StreamID: {%s}, Version: {%d}, EventType: {%s}, AggregateType: {%s}, Metadata: {%s}, TimeStamp: {%s}",
		e.StreamID,
		e.Version,
		e.EventType,
		e.AggregateType,
		string(e.Metadata),
		e.Timestamp.String(),
	)
}
