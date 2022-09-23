package serializer

import "reflect"

type EventSerializer interface {
	Serialize(event interface{}) (*EventSerializationResult, error)
	Deserialize(data []byte, eventType string, contentType string) (interface{}, error)
	DeserializeType(data []byte, eventType reflect.Type, contentType string) (interface{}, error)
	DeserializeMessage(data []byte, eventType string, contentType string) (interface{}, error)
	DeserializeEvent(data []byte, eventType string, contentType string) (interface{}, error)
	ContentType() string
}

type EventSerializationResult struct {
	Data        []byte
	ContentType string
}
