package messaging

// https://opentelemetry.io/docs/specs/otel/trace/semantic_conventions/messaging/
const (
	System          = "messaging.system"
	Destination     = "messaging.destination"
	DestinationKind = "messaging.destination_kind"
	Url             = "messaging.url"
	MessageId       = "messaging.message_id"
	ConversationId  = "messaging.conversation_id"
	CorrelationId   = "messaging.correlation_id"
	CausationId     = "messaging.causation_id"
	Operation       = "messaging.operation"
)
