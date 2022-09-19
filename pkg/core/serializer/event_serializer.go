package serializer

type EventSerializer interface {
	Serialize(event interface{}) (*EventSerializationResult, error)
	Deserialize(data []byte, eventType string, contentType string) (interface{}, error)
	ContentType() string
}

type EventSerializationResult struct {
	Data        []byte
	ContentType string
	EventType   string
}
