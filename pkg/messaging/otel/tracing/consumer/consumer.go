package consumer

import (
	"context"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/metadata"
	messageHeader "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/message_header"
	messageTracing "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/otel/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	tracingHeaders "github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing/tracing_headers"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"time"
)

//https://devandchill.com/posts/2021/12/go-step-by-step-guide-for-implementing-tracing-on-a-microservices-architecture-2/2/
//https://github.com/open-telemetry/opentelemetry-go-contrib/blob/e84d6d6575e3c3eabcf3204ac88550258673ed3c/instrumentation/github.com/Shopify/sarama/otelsarama/dispatcher.go
//https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/messaging/
//https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#messaging-attributes
//https://opentelemetry.io/docs/instrumentation/go/manual/#semantic-attributes
//https://trstringer.com/otel-part5-propagation/

func StartConsumerSpan(ctx context.Context, meta *metadata.Metadata, payload string, consumerTracingOptions *ConsumerTracingOptions) (context.Context, trace.Span) {
	ctx = addAfterBaggage(ctx, meta)

	// If there's a span context in the message, use that as the parent context.
	// extracts the tracing from the header and puts it into the context
	carrier := messageTracing.NewMessageCarrier(meta)
	parentSpanContext := otel.GetTextMapPropagator().Extract(ctx, carrier)

	opts := getTraceOptions(meta, payload, consumerTracingOptions)

	//https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#span-name
	// SpanName = Destination Name + Operation Name
	ctx, span := messageTracing.MessagingTracer.Start(parentSpanContext, fmt.Sprintf("%s %s", consumerTracingOptions.Destination, "receive"), opts...)

	span.AddEvent(fmt.Sprintf("start consuming message '%s' from the broker", messageHeader.GetMessageName(*meta)))

	// Emulate Work loads
	time.Sleep(1 * time.Second)

	// we don't want next trace (AfterConsume) becomes child of this span, so we should not use new ctx for (AfterConsume) span. if already exists a span on ctx next span will be a child of that span
	return ctx, span
}

func FinishConsumerSpan(span trace.Span, err error) error {
	messageName := tracing.GetSpanAttribute(span, messageTracing.MessageName).Value.AsString()

	if err != nil {
		span.AddEvent(fmt.Sprintf("failed to consume message '%s' from the broker", messageName))
		_ = messageTracing.TraceMessagingErrFromSpan(span, err)
	}

	span.SetAttributes(
		attribute.Key(tracing.SpanId).String(span.SpanContext().SpanID().String()), // current span id
	)

	span.AddEvent(fmt.Sprintf("message '%s' consumed from the broker succesfully", messageName))
	span.End()

	return err
}

func getTraceOptions(meta *metadata.Metadata, payload string, consumerTracingOptions *ConsumerTracingOptions) []trace.SpanStartOption {
	correlationId := messageHeader.GetCorrelationId(*meta)

	//https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#topic-with-multiple-consumers
	//https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/semantic_conventions/messaging.md#batch-receiving
	attrs := []attribute.KeyValue{
		semconv.MessageIDKey.String(messageHeader.GetMessageId(*meta)),
		semconv.MessagingConversationIDKey.String(correlationId),
		semconv.MessagingOperationReceive,
		attribute.Key(tracing.TraceId).String(tracingHeaders.GetTracingTraceId(*meta)),
		attribute.Key(tracing.Traceparent).String(tracingHeaders.GetTracingTraceparent(*meta)),
		attribute.Key(tracing.ParentSpanId).String(tracingHeaders.GetTracingParentSpanId(*meta)),
		attribute.Key(tracing.Timestamp).Int64(time.Now().UnixMilli()),
		attribute.Key(messageTracing.MessageType).String(messageHeader.GetMessageType(*meta)),
		attribute.Key(messageTracing.MessageName).String(messageHeader.GetMessageName(*meta)),
		attribute.Key(messageTracing.Payload).String(payload),
		attribute.String(messageTracing.Headers, meta.ToJson()),
		semconv.MessagingDestinationKey.String(consumerTracingOptions.Destination),
		semconv.MessagingSystemKey.String(consumerTracingOptions.MessagingSystem),
		semconv.MessagingDestinationKindKey.String(consumerTracingOptions.DestinationKind),
	}

	if consumerTracingOptions.OtherAttributes != nil && len(consumerTracingOptions.OtherAttributes) > 0 {
		attrs = append(attrs, consumerTracingOptions.OtherAttributes...)
	}

	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindConsumer),
	}
	return opts
}

func addAfterBaggage(ctx context.Context, meta *metadata.Metadata) context.Context {
	correlationId := messageHeader.GetCorrelationId(*meta)

	correlationIdBag, _ := baggage.NewMember(string(semconv.MessagingConversationIDKey), correlationId)
	messageIdBag, _ := baggage.NewMember(string(semconv.MessageIDKey), messageHeader.GetMessageId(*meta))
	b, _ := baggage.New(correlationIdBag, messageIdBag)
	ctx = baggage.ContextWithBaggage(ctx, b)

	// new context including baggage
	return ctx
}
