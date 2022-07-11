package es

import (
	"encoding/json"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/esdb"
	uuid "github.com/satori/go.uuid"
	"time"
)

// EventType is the type of any event, used as its unique identifier.
type EventType string

// Event is an internal representation of an event, returned when the Aggregate
// uses NewEvent to create a new event. The events loaded from the db is
// represented by each DBs internal event type, implementing Event.
type Event struct {
	EventID       string
	EventType     string
	Data          []byte
	Timestamp     time.Time
	AggregateType AggregateType
	AggregateID   string
	Version       int64
	Metadata      []byte
}

// NewBaseEvent new base Event constructor with configured EventID, Aggregate properties and Timestamp.
func NewBaseEvent(aggregate Aggregate, eventType string) Event {
	return Event{
		EventID:       uuid.NewV4().String(),
		AggregateType: aggregate.GetType(),
		AggregateID:   aggregate.GetID(),
		Version:       aggregate.GetVersion(),
		EventType:     eventType,
		Timestamp:     time.Now(),
	}
}

func NewEventFromRecorded(event *esdb.RecordedEvent) Event {
	return Event{
		EventID:     event.EventID.String(),
		EventType:   event.EventType,
		Data:        event.Data,
		Timestamp:   event.CreatedDate,
		AggregateID: event.StreamID,
		Version:     int64(event.EventNumber),
		Metadata:    event.UserMetadata,
	}
}

func NewEventFromEventData(event esdb.EventData) Event {
	return Event{
		EventID:   event.EventID.String(),
		EventType: event.EventType,
		Data:      event.Data,
		Metadata:  event.Metadata,
	}
}

func EventFromEventData(recordedEvent esdb.RecordedEvent) (Event, error) {
	return Event{
		EventID:     recordedEvent.EventID.String(),
		EventType:   recordedEvent.EventType,
		Data:        recordedEvent.Data,
		Timestamp:   recordedEvent.CreatedDate,
		AggregateID: recordedEvent.StreamID,
		Version:     int64(recordedEvent.Position.Commit),
		Metadata:    nil,
	}, nil
}

func (e *Event) ToEventData() esdb.EventData {
	return esdb.EventData{
		EventType:   e.EventType,
		ContentType: esdb.JsonContentType,
		Data:        e.Data,
		Metadata:    e.Metadata,
	}
}

// GetEventID get EventID of the Event.
func (e *Event) GetEventID() string {
	return e.EventID
}

// GetTimeStamp get timestamp of the Event.
func (e *Event) GetTimeStamp() time.Time {
	return e.Timestamp
}

// GetData The data attached to the Event serialized to bytes.
func (e *Event) GetData() []byte {
	return e.Data
}

// SetData add the data attached to the Event serialized to bytes.
func (e *Event) SetData(data []byte) *Event {
	e.Data = data
	return e
}

// GetJsonData json unmarshal data attached to the Event.
func (e *Event) GetJsonData(data interface{}) error {
	return json.Unmarshal(e.GetData(), data)
}

// SetJsonData serialize to json and set data attached to the Event.
func (e *Event) SetJsonData(data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	e.Data = dataBytes
	return nil
}

// GetEventType returns the EventType of the event.
func (e *Event) GetEventType() string {
	return e.EventType
}

// GetAggregateType is the AggregateType that the Event can be applied to.
func (e *Event) GetAggregateType() AggregateType {
	return e.AggregateType
}

// SetAggregateType set the AggregateType that the Event can be applied to.
func (e *Event) SetAggregateType(aggregateType AggregateType) {
	e.AggregateType = aggregateType
}

// GetAggregateID is the ID of the Aggregate that the Event belongs to
func (e *Event) GetAggregateID() string {
	return e.AggregateID
}

// GetVersion is the version of the Aggregate after the Event has been applied.
func (e *Event) GetVersion() int64 {
	return e.Version
}

// SetVersion set the version of the Aggregate.
func (e *Event) SetVersion(aggregateVersion int64) {
	e.Version = aggregateVersion
}

// GetMetadata is app-specific metadata such as request ID, originating user etc.
func (e *Event) GetMetadata() []byte {
	return e.Metadata
}

// SetMetadata add app-specific metadata serialized as json for the Event.
func (e *Event) SetMetadata(metaData interface{}) error {

	metaDataBytes, err := json.Marshal(metaData)
	if err != nil {
		return err
	}

	e.Metadata = metaDataBytes
	return nil
}

// GetJsonMetadata unmarshal app-specific metadata serialized as json for the Event.
func (e *Event) GetJsonMetadata(metaData interface{}) error {
	return json.Unmarshal(e.GetMetadata(), metaData)
}

// GetString A string representation of the Event.
func (e *Event) GetString() string {
	return fmt.Sprintf("event: %+v", e)
}

func (e *Event) String() string {
	return fmt.Sprintf("(Event): AggregateID: {%s}, Version: {%d}, EventType: {%s}, AggregateType: {%s}, Metadata: {%s}, TimeStamp: {%s}",
		e.AggregateID,
		e.Version,
		e.EventType,
		e.AggregateType,
		string(e.Metadata),
		e.Timestamp.String(),
	)
}
