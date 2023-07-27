package tracing

import (
    "context"

    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/baggage"
    "go.opentelemetry.io/otel/trace"
    "google.golang.org/grpc/metadata"

    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
)

var GrpcTracer trace.Tracer

//https://opentelemetry.io/docs/instrumentation/go/manual/#semantic-attributes
//https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/rpc/

func init() {
	//
	GrpcTracer = tracing.NewCustomTracer("github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc") //instrumentation name
}

// StartGrpcServerTracerSpan uses when grpc otel middleware is off and create a span on 'grpc' tracer
func StartGrpcServerTracerSpan(ctx context.Context, operationName string) (context context.Context, span trace.Span, deferSpan func()) {
	requestMetadata, _ := metadata.FromIncomingContext(ctx)
	metadataCopy := requestMetadata.Copy()

	bags, spanCtx := otelgrpc.Extract(ctx, &metadataCopy)
	ctx = baggage.ContextWithBaggage(ctx, bags)

	attrs := []attribute.KeyValue{otelgrpc.RPCSystemGRPC}

	context, span = GrpcTracer.Start(
		trace.ContextWithRemoteSpanContext(ctx, spanCtx),
		operationName,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attrs...),
	)

	return context, span, func() {
		span.End()
	}
}
