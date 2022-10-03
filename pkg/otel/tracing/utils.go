package tracing

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/metadata"
	errorUtils "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils/error_utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	trace2 "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"reflect"
)

type traceContextKeyType int

const parentSpanKey traceContextKeyType = iota + 1

func TraceErrFromContext(ctx context.Context, err error) error {
	//https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
	span := trace.SpanFromContext(ctx)
	defer span.End()

	if err != nil {
		stackTraceError := errorUtils.ErrorsWithStack(err)
		span.SetStatus(codes.Error, "")
		span.SetAttributes(attribute.String(ErrorMessage, stackTraceError))
		span.RecordError(err)
	}

	return err
}

func TraceErrFromSpan(span trace.Span, err error) error {
	if err != nil {
		stackTraceError := errorUtils.ErrorsWithStack(err)
		span.SetStatus(codes.Error, "")
		span.SetAttributes(attribute.String(ErrorMessage, stackTraceError))
		span.RecordError(err)
	}

	return err
}

func GetParentSpanContext(span trace.Span) trace.SpanContext {
	readWriteSpan, ok := span.(trace2.ReadWriteSpan)
	if !ok {
		return *new(trace.SpanContext)
	}

	return readWriteSpan.Parent()
}

func ContextWithParentSpan(parent context.Context, span trace.Span) context.Context {
	return context.WithValue(parent, parentSpanKey, span)
}

func ParentSpanFromContext(ctx context.Context) trace.Span {
	_, nopSpan := trace.NewNoopTracerProvider().Tracer("").Start(ctx, "")
	if ctx == nil {
		return nopSpan
	}
	if span, ok := ctx.Value(parentSpanKey).(trace.Span); ok {
		return span
	}
	return nopSpan
}

func CopyFromParentSpanAttribute(ctx context.Context, span trace.Span, attributeName string, parentAttributeName string) {
	parentAtt := GetParentSpanAttribute(ctx, parentAttributeName)
	if reflect.ValueOf(parentAtt).IsZero() {
		return
	}
	span.SetAttributes(attribute.String(attributeName, parentAtt.Value.AsString()))
}

func CopyFromParentSpanAttributeIfNotSet(ctx context.Context, span trace.Span, attributeName string, attributeValue string, parentAttributeName string) {
	if attributeValue != "" {
		span.SetAttributes(attribute.String(attributeName, attributeValue))
		return
	}
	CopyFromParentSpanAttribute(ctx, span, attributeName, parentAttributeName)
}

func GetParentSpanAttribute(ctx context.Context, parentAttributeName string) attribute.KeyValue {
	parentSpan := ParentSpanFromContext(ctx)
	readWriteSpan, ok := parentSpan.(trace2.ReadWriteSpan)
	if !ok {
		return *new(attribute.KeyValue)
	}
	att := linq.From(readWriteSpan.Attributes()).FirstWithT(func(att attribute.KeyValue) bool { return string(att.Key) == parentAttributeName })

	return att.(attribute.KeyValue)
}

func GetSpanAttributeFromCurrentContext(ctx context.Context, attributeName string) attribute.KeyValue {
	span := trace.SpanFromContext(ctx)
	readWriteSpan, ok := span.(trace2.ReadWriteSpan)
	if !ok {
		return *new(attribute.KeyValue)
	}
	att := linq.From(readWriteSpan.Attributes()).FirstWithT(func(att attribute.KeyValue) bool { return string(att.Key) == attributeName })

	return att.(attribute.KeyValue)
}

func GetSpanAttribute(span trace.Span, attributeName string) attribute.KeyValue {
	readWriteSpan, ok := span.(trace2.ReadWriteSpan)
	if !ok {
		return *new(attribute.KeyValue)
	}
	att := linq.From(readWriteSpan.Attributes()).FirstWithT(func(att attribute.KeyValue) bool { return string(att.Key) == attributeName })

	return att.(attribute.KeyValue)
}

func MapsToAttributes(maps map[string]interface{}) []attribute.KeyValue {
	var att []attribute.KeyValue
	for key, val := range maps {
		switch val.(type) {
		case string:
			att = append(att, attribute.String(key, val.(string)))
		case int64:
			att = append(att, attribute.Int64(key, val.(int64)))
		case int, int32:
			att = append(att, attribute.Int(key, val.(int)))
		case float64, float32:
			att = append(att, attribute.Float64(key, val.(float64)))
		case bool:
			att = append(att, attribute.Bool(key, val.(bool)))
		}
	}

	return att
}

func MetadataToSet(meta metadata.Metadata) attribute.Set {
	var keyValue []attribute.KeyValue
	for key, val := range meta {
		keyValue = append(keyValue, attribute.String(key, val.(string)))
	}
	return attribute.NewSet(keyValue...)
}
