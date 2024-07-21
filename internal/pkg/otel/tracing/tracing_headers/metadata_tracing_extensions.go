package tracingHeaders

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/metadata"
)

func GetTracingTraceId(m metadata.Metadata) string {
	return m.GetString(TraceId)
}

func GetTracingParentSpanId(m metadata.Metadata) string {
	return m.GetString(ParentSpanId)
}

func GetTracingTraceparent(m metadata.Metadata) string {
	return m.GetString(Traceparent)
}
