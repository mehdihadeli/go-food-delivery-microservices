package consumer

import "go.opentelemetry.io/otel/attribute"

type ConsumerTracingOptions struct {
	MessagingSystem string
	DestinationKind string
	Destination     string
	OtherAttributes []attribute.KeyValue
}
