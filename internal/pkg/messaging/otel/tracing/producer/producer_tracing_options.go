package producer

import "go.opentelemetry.io/otel/attribute"

type ProducerTracingOptions struct {
	MessagingSystem string
	DestinationKind string
	Destination     string
	OtherAttributes []attribute.KeyValue
}
