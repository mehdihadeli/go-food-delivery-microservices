package producer

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/metadata"
	messageHeader "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/message_header"
	messageTracing "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"time"
)

//https://devandchill.com/posts/2021/12/go-step-by-step-guide-for-implementing-tracing-on-a-microservices-architecture-2/2/
//https://github.com/open-telemetry/opentelemetry-go-contrib/blob/v0.12.0/instrumentation/github.com/Shopify/sarama/otelsarama/producer.go
//https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/messaging/
//https://opentelemetry.io/docs/instrumentation/go/manual/#semantic-attributes
//https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#messaging-attributes
//https://trstringer.com/otel-part5-propagation/

func StartProducerSpan(ctx context.Context, message types.IMessage, meta *metadata.Metadata, payload string, producerTracingOptions *ProducerTracingOptions) (context.Context, trace.Span) {
	ctx = addAfterBaggage(ctx, message, meta)

	// If there's a span context in the message, use that as the parent context.
	// extracts the tracing from the header and puts it into the context
	carrier := messageTracing.NewMessageCarrier(meta)
	parentSpanContext := otel.GetTextMapPropagator().Extract(ctx, carrier)

	opts := getTraceOptions(meta, message, payload, producerTracingOptions)

	//https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#span-name
	// SpanName = Destination Name + Operation Name
	ctx, span := messageTracing.MessagingTracer.Start(parentSpanContext, fmt.Sprintf("%s %s", producerTracingOptions.Destination, "send"), opts...)

	span.AddEvent(fmt.Sprintf("start publishing message '%s' to the broker", messageHeader.GetMessageName(*meta)))

	// Injects current span context, so consumers can use it to propagate span.
	// injects the tracing from the context into the header map
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	// we don't want next trace (AfterProduce) becomes child of this span, so we should not use new ctx for (AfterProducer) span. if already exists a span on ctx next span will be a child of that span
	return ctx, span
}

func FinishProducerSpan(span trace.Span, err error) error {
	messageName := tracing.GetSpanAttribute(span, messageTracing.MessageName).Value.AsString()

	if err != nil {
		span.AddEvent(fmt.Sprintf("failed to publsih message '%s' to the broker", messageName))
		_ = tracing.TraceErrFromSpan(span, err)
	}
	span.SetAttributes(
		attribute.Key(tracing.TraceId).String(span.SpanContext().TraceID().String()),
		attribute.Key(tracing.SpanId).String(span.SpanContext().SpanID().String()), // current span id
	)

	span.AddEvent(fmt.Sprintf("message '%s' published to the broker succesfully", messageName))
	span.End()

	return err
}

func getTraceOptions(meta *metadata.Metadata, message types.IMessage, payload string, producerTracingOptions *ProducerTracingOptions) []trace.SpanStartOption {
	correlationId := messageHeader.GetCorrelationId(*meta)

	//https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#topic-with-multiple-consumers
	//https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#batch-receiving
	attrs := []attribute.KeyValue{
		semconv.MessageIDKey.String(message.GeMessageId()),
		semconv.MessagingConversationIDKey.String(correlationId),
		attribute.Key(messageTracing.MessageType).String(message.GetEventTypeName()),
		attribute.Key(messageTracing.MessageName).String(messageHeader.GetMessageName(*meta)),
		attribute.Key(messageTracing.Payload).String(payload),
		attribute.String(messageTracing.Headers, meta.ToJson()),
		attribute.Key(tracing.Timestamp).Int64(time.Now().UnixMilli()),
		semconv.MessagingDestinationKey.String(producerTracingOptions.Destination),
		semconv.MessagingSystemKey.String(producerTracingOptions.MessagingSystem),
		semconv.MessagingDestinationKindKey.String(producerTracingOptions.DestinationKind),
		semconv.MessagingOperationKey.String("send"),
	}

	if producerTracingOptions.OtherAttributes != nil && len(producerTracingOptions.OtherAttributes) > 0 {
		attrs = append(attrs, producerTracingOptions.OtherAttributes...)
	}

	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindProducer),
	}
	return opts
}

func addAfterBaggage(ctx context.Context, message types.IMessage, meta *metadata.Metadata) context.Context {
	correlationId := messageHeader.GetCorrelationId(*meta)

	correlationIdBag, _ := baggage.NewMember(string(semconv.MessagingConversationIDKey), correlationId)
	messageIdBag, _ := baggage.NewMember(string(semconv.MessageIDKey), message.GeMessageId())
	b, _ := baggage.New(correlationIdBag, messageIdBag)
	ctx = baggage.ContextWithBaggage(ctx, b)

	// new context including baggage
	return ctx
}
