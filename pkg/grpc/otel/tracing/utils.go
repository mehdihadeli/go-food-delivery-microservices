package tracing

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/grpc/grpcErrors"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	errorUtils "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils/error_utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// TraceGrpcErrFromSpan setting span with status error with error message
func TraceGrpcErrFromSpan(span trace.Span, err error) error {
	if err != nil {
		stackTraceError := errorUtils.ErrorsWithStack(err)
		span.SetStatus(codes.Error, "")
		span.SetAttributes(attribute.String(GrpcErrorMessage, stackTraceError))
		if customErrors.IsCustomError(err) {
			grpcErr := grpcErrors.ParseError(err)
			span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(grpcErr.GetStatus())))
		}
		span.RecordError(err)
	}

	return err
}

// TraceGrpcErrFromSpanWithCode setting span with status error with error message
func TraceGrpcErrFromSpanWithCode(span trace.Span, err error, code int) error {
	if err != nil {
		stackTraceError := errorUtils.ErrorsWithStack(err)
		span.SetStatus(codes.Error, "")
		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(code))
		span.SetAttributes(attribute.String(GrpcErrorMessage, stackTraceError))
		span.RecordError(err)
	}

	return err
}
