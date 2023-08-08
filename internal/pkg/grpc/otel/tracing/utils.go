package tracing

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/grpc/grpcErrors"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	errorUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/error_utils"
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
