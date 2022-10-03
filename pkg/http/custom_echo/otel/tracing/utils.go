package tracing

import (
	"context"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/problemDetails"
	errorUtils "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils/error_utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// TraceHttpErrFromSpan setting span with status error with error message
func TraceHttpErrFromSpan(span trace.Span, err error) error {
	if err != nil {
		stackTraceError := errorUtils.ErrorsWithStack(err)
		span.SetStatus(codes.Error, "")
		span.SetAttributes(attribute.String(HttpErrorMessage, stackTraceError))
		if customErrors.IsCustomError(err) {
			httpError := problemDetails.ParseError(err)
			span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(httpError.GetStatus())...)
		}
		span.RecordError(err)
	}

	return err
}

// TraceHttpErrFromSpanWithCode setting span with status error with error message
func TraceHttpErrFromSpanWithCode(span trace.Span, err error, code int) error {
	if err != nil {
		stackTraceError := errorUtils.ErrorsWithStack(err)
		span.SetStatus(codes.Error, "")
		span.SetAttributes(semconv.HTTPAttributesFromHTTPStatusCode(code)...)
		span.SetAttributes(attribute.String(HttpErrorMessage, stackTraceError))
		span.RecordError(err)
	}

	return err
}

func TraceHttpErrFromContext(ctx context.Context, err error) error {
	//https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
	span := trace.SpanFromContext(ctx)
	defer span.End()

	if err != nil {
		stackTraceError := errorUtils.ErrorsWithStack(err)
		span.SetStatus(codes.Error, "")
		span.SetAttributes(attribute.String(HttpErrorMessage, stackTraceError))
		span.RecordError(err)
	}

	return err
}
