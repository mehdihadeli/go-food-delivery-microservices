package tracing

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/otel/tracing"
	"go.opentelemetry.io/otel/trace"
)

var MessagingTracer trace.Tracer

func init() {
	MessagingTracer = tracing.NewCustomTracer("github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging") //instrumentation name
}
