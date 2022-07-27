package domain

type EventEnvelope struct {
	EventData interface{}
	Metadata  *Metadata
}
