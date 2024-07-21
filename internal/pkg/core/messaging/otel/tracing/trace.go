package tracing

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/tracing"

	"go.opentelemetry.io/otel/trace"
)

var MessagingTracer trace.Tracer

func init() {
	MessagingTracer = tracing.NewAppTracer(
		"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/messaging",
	) // instrumentation name
}
