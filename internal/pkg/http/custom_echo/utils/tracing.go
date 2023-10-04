package utils

import (
	"context"
	"net/http"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/custom_echo/constants"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/problemDetails"
	errorUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/error_utils"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// HttpTraceFromSpan create an error span if we have an error and a successful span when error is nil
func HttpTraceFromSpan(span trace.Span, err error) error {
	isError := err != nil

	if customErrors.IsCustomError(err) {
		httpError := problemDetails.ParseError(err)

		return HttpTraceFromSpanWithCode(
			span,
			err,
			httpError.GetStatus(),
		)
	}

	var (
		status int
		code   codes.Code
	)

	if isError {
		status = http.StatusInternalServerError
		code = codes.Error
	} else {
		status = http.StatusOK
		code = codes.Ok
	}

	span.SetStatus(code, "")
	span.SetAttributes(
		semconv.HTTPStatusCode(status),
	)

	if isError {
		stackTraceError := errorUtils.ErrorsWithStack(err)

		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(constants.Otel.HttpErrorMessage, stackTraceError),
		)
		span.RecordError(err)
	}

	return err
}

// HttpTraceFromSpanWithCode create an error span with specific status code if we have an error and a successful span when error is nil with a specific status
func HttpTraceFromSpanWithCode(span trace.Span, err error, code int) error {
	if err != nil {
		stackTraceError := errorUtils.ErrorsWithStack(err)

		// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
		span.SetAttributes(
			attribute.String(constants.Otel.HttpErrorMessage, stackTraceError),
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

// HttpTraceFromContext create an error span if we have an error and a successful span when error is nil
func HttpTraceFromContext(ctx context.Context, err error) error {
	// https://opentelemetry.io/docs/instrumentation/go/manual/#record-errors
	span := trace.SpanFromContext(ctx)

	defer span.End()

	return HttpTraceFromSpan(span, err)
}
