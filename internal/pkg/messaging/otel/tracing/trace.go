package tracing

import (
    "go.opentelemetry.io/otel/trace"

    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
)

var MessagingTracer trace.Tracer

func init() {
	MessagingTracer = tracing.NewCustomTracer("github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging") //instrumentation name
}
