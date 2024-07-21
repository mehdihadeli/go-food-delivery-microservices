package grpcerrors

import (
	"context"
	"database/sql"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/constants"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	errorUtils "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils/errorutils"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"google.golang.org/grpc/codes"
)

// https://github.com/grpc/grpc/blob/master/doc/http-grpc-status-mapping.md
// https://github.com/grpc/grpc/blob/master/doc/statuscodes.md

func ParseError(err error) GrpcErr {
	customErr := customErrors.GetCustomError(err)
	var validatorErr validator.ValidationErrors
	stackTrace := errorUtils.ErrorsWithStack(err)

	if err != nil && customErr != nil {
		switch {
		case customErrors.IsDomainError(err, customErr.Status()):
			return NewDomainGrpcError(
				codes.Code(customErr.Status()),
				customErr.Error(),
				stackTrace,
			)
		case customErrors.IsApplicationError(err, customErr.Status()):
			return NewApplicationGrpcError(
				codes.Code(customErr.Status()),
				customErr.Error(),
				stackTrace,
			)
		case customErrors.IsApiError(err, customErr.Status()):
			return NewApiGrpcError(
				codes.Code(customErr.Status()),
				customErr.Error(),
				stackTrace,
			)
		case customErrors.IsBadRequestError(err):
			return NewBadRequestGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsNotFoundError(err):
			return NewNotFoundErrorGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsValidationError(err):
			return NewValidationGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsUnAuthorizedError(err):
			return NewUnAuthorizedErrorGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsForbiddenError(err):
			return NewForbiddenGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsConflictError(err):
			return NewConflictGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsInternalServerError(err):
			return NewInternalServerGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsCustomError(err):
			return NewGrpcError(
				codes.Code(customErr.Status()),
				codes.Code(customErr.Status()).String(),
				customErr.Error(),
				stackTrace,
			)
		case customErrors.IsUnMarshalingError(err):
			return NewInternalServerGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsMarshalingError(err):
			return NewInternalServerGrpcError(customErr.Error(), stackTrace)
		default:
			return NewInternalServerGrpcError(err.Error(), stackTrace)
		}
	} else if err != nil && customErr == nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return NewNotFoundErrorGrpcError(err.Error(), stackTrace)
		case errors.Is(err, context.DeadlineExceeded):
			return NewGrpcError(
				codes.DeadlineExceeded,
				constants.ErrRequestTimeoutTitle,
				err.Error(),
				stackTrace,
			)
		case errors.As(err, &validatorErr):
			return NewValidationGrpcError(validatorErr.Error(), stackTrace)
		}
	}

	return nil
}
