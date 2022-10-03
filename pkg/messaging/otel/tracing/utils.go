package tracing

import (
	"context"
	errorUtils "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils/error_utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TraceMessagingErrFromSpan setting span with status error with error message
func TraceMessagingErrFromSpan(span trace.Span, err error) error {
	if err != nil {
		stackTraceError := errorUtils.ErrorsWithStack(err)
		span.SetStatus(codes.Error, "")
		span.SetAttributes(attribute.String(MessagingErrorMessage, stackTraceError))
		span.RecordError(err)
	}

	return err
}

func TraceMessagingErrFromContext(ctx context.Context, err error) error {
	//https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
	span := trace.SpanFromContext(ctx)
	defer span.End()

	if err != nil {
		stackTraceError := errorUtils.ErrorsWithStack(err)
		span.SetStatus(codes.Error, "")
		span.SetAttributes(attribute.String(MessagingErrorMessage, stackTraceError))
		span.RecordError(err)
	}

	return err
}
