package utils

import (
	"context"
	"net/http"
	"reflect"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/metadata"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/grpc/grpcerrors"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	problemdetails "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/problemdetails"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/otel/constants/telemetrytags"
	errorUtils "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils/errorutils"

	"github.com/ahmetb/go-linq/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	trace2 "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

type traceContextKeyType int

const parentSpanKey traceContextKeyType = iota + 1

// HttpTraceStatusFromSpan create an error span if we have an error and a successful span when error is nil
func HttpTraceStatusFromSpan(span trace.Span, err error) error {
	isError := err != nil

	if customErrors.IsCustomError(err) {
		httpError := problemdetails.ParseError(err)

		return HttpTraceStatusFromSpanWithCode(
			span,
			err,
			httpError.GetStatus(),
		)
	}

	var (
		status      int
		code        codes.Code
		description = ""
	)

	if isError {
		status = http.StatusInternalServerError
		code = codes.Error
		description = err.Error()
	} else {
		status = http.StatusOK
		code = codes.Ok
	}

	span.SetStatus(code, description)
	span.SetAttributes(
		semconv.HTTPStatusCode(status),
	)

	if isError {
		stackTraceError := errorUtils.ErrorsWithStack(err)

		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(telemetrytags.Exceptions.Message, err.Error()),
			attribute.String(telemetrytags.Exceptions.Stacktrace, stackTraceError),
		)
		span.RecordError(err)
	}

	return err
}

func TraceStatusFromSpan(span trace.Span, err error) error {
	isError := err != nil

	var (
		code        codes.Code
		description = ""
	)

	if isError {
		code = codes.Error
		description = err.Error()
	} else {
		code = codes.Ok
	}

	span.SetStatus(code, description)

	if isError {
		stackTraceError := errorUtils.ErrorsWithStack(err)

		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(telemetrytags.Exceptions.Message, err.Error()),
			attribute.String(telemetrytags.Exceptions.Stacktrace, stackTraceError),
		)
		span.RecordError(err)
	}

	return err
}

func TraceErrStatusFromSpan(span trace.Span, err error) error {
	isError := err != nil

	span.SetStatus(codes.Error, err.Error())

	if isError {
		stackTraceError := errorUtils.ErrorsWithStack(err)

		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(telemetrytags.Exceptions.Message, err.Error()),
			attribute.String(telemetrytags.Exceptions.Stacktrace, stackTraceError),
		)
		span.RecordError(err)
	}

	return err
}

// HttpTraceStatusFromSpanWithCode create an error span with specific status code if we have an error and a successful span when error is nil with a specific status
func HttpTraceStatusFromSpanWithCode(
	span trace.Span,
	err error,
	code int,
) error {
	if err != nil {
		stackTraceError := errorUtils.ErrorsWithStack(err)

		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(telemetrytags.Exceptions.Message, err.Error()),
			attribute.String(telemetrytags.Exceptions.Stacktrace, stackTraceError),
		)
		span.RecordError(err)
	}

	if code > 0 {
		// httpconv doesn't exist in semconv v1.21.0, and it moved to `opentelemetry-go-contrib` pkg
		// https://github.com/open-telemetry/opentelemetry-go/pull/4362
		// https://github.com/open-telemetry/opentelemetry-go/issues/4081
		// using ClientStatus instead of ServerStatus for consideration of 4xx status as error
		span.SetStatus(httpconv.ClientStatus(code))
		span.SetAttributes(semconv.HTTPStatusCode(code))
	} else {
		span.SetStatus(codes.Error, "")
		span.SetAttributes(semconv.HTTPStatusCode(http.StatusInternalServerError))
	}

	return err
}

// HttpTraceStatusFromContext create an error span if we have an error and a successful span when error is nil
func HttpTraceStatusFromContext(ctx context.Context, err error) error {
	// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
	span := trace.SpanFromContext(ctx)

	defer span.End()

	return HttpTraceStatusFromSpan(span, err)
}

func TraceStatusFromContext(ctx context.Context, err error) error {
	// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
	span := trace.SpanFromContext(ctx)

	defer span.End()

	return TraceStatusFromSpan(span, err)
}

func TraceErrStatusFromContext(ctx context.Context, err error) error {
	// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
	span := trace.SpanFromContext(ctx)

	defer span.End()

	return TraceErrStatusFromSpan(span, err)
}

// GrpcTraceErrFromSpan setting span with status error with error message
func GrpcTraceErrFromSpan(span trace.Span, err error) error {
	isError := err != nil

	span.SetStatus(codes.Error, err.Error())

	if isError {
		stackTraceError := errorUtils.ErrorsWithStack(err)
		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(telemetrytags.Exceptions.Message, err.Error()),
			attribute.String(telemetrytags.Exceptions.Stacktrace, stackTraceError),
		)

		if customErrors.IsCustomError(err) {
			grpcErr := grpcerrors.ParseError(err)
			span.SetAttributes(
				semconv.RPCGRPCStatusCodeKey.Int(int(grpcErr.GetStatus())),
			)
		}

		span.RecordError(err)
	}

	return err
}

// GrpcTraceErrFromSpanWithCode setting span with status error with error message
func GrpcTraceErrFromSpanWithCode(span trace.Span, err error, code int) error {
	isError := err != nil

	span.SetStatus(codes.Error, err.Error())

	if isError {
		stackTraceError := errorUtils.ErrorsWithStack(err)
		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(telemetrytags.Exceptions.Message, err.Error()),
			attribute.String(telemetrytags.Exceptions.Stacktrace, stackTraceError),
		)
		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(code))
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

func ContextWithParentSpan(
	parent context.Context,
	span trace.Span,
) context.Context {
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

func CopyFromParentSpanAttribute(
	ctx context.Context,
	span trace.Span,
	attributeName string,
	parentAttributeName string,
) {
	parentAtt := GetParentSpanAttribute(ctx, parentAttributeName)
	if reflect.ValueOf(parentAtt).IsZero() {
		return
	}

	span.SetAttributes(
		attribute.String(attributeName, parentAtt.Value.AsString()),
	)
}

func CopyFromParentSpanAttributeIfNotSet(
	ctx context.Context,
	span trace.Span,
	attributeName string,
	attributeValue string,
	parentAttributeName string,
) {
	if attributeValue != "" {
		span.SetAttributes(attribute.String(attributeName, attributeValue))
		return
	}
	CopyFromParentSpanAttribute(ctx, span, attributeName, parentAttributeName)
}

func GetParentSpanAttribute(
	ctx context.Context,
	parentAttributeName string,
) attribute.KeyValue {
	parentSpan := ParentSpanFromContext(ctx)
	readWriteSpan, ok := parentSpan.(trace2.ReadWriteSpan)
	if !ok {
		return *new(attribute.KeyValue)
	}
	att := linq.From(readWriteSpan.Attributes()).
		FirstWithT(func(att attribute.KeyValue) bool { return string(att.Key) == parentAttributeName })

	return att.(attribute.KeyValue)
}

func GetSpanAttributeFromCurrentContext(
	ctx context.Context,
	attributeName string,
) attribute.KeyValue {
	span := trace.SpanFromContext(ctx)
	readWriteSpan, ok := span.(trace2.ReadWriteSpan)
	if !ok {
		return *new(attribute.KeyValue)
	}
	att := linq.From(readWriteSpan.Attributes()).
		FirstWithT(func(att attribute.KeyValue) bool { return string(att.Key) == attributeName })

	return att.(attribute.KeyValue)
}

func GetSpanAttribute(
	span trace.Span,
	attributeName string,
) attribute.KeyValue {
	readWriteSpan, ok := span.(trace2.ReadWriteSpan)
	if !ok {
		return *new(attribute.KeyValue)
	}

	att := linq.From(readWriteSpan.Attributes()).
		FirstWithT(func(att attribute.KeyValue) bool { return string(att.Key) == attributeName })

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
