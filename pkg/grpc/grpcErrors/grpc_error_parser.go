package grpcErrors

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger/defaultLogger"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//ErrGrpcResponse get gRPC error response
func ErrGrpcResponse(err error) error {
	grpcErr := ParseError(err)
	stackTrace := grpcErr.GetStackTrace()
	defaultLogger.Logger.Error(err.Error())

	if core.IsDevelopment() {
		stackTrace := stackTrace
		fmt.Println(stackTrace)
	}

	return status.Error(grpcErr.GetStatus(), grpcErr.ToJson())
}

//https://github.com/grpc/grpc/blob/master/doc/http-grpc-status-mapping.md
//https://github.com/grpc/grpc/blob/master/doc/statuscodes.md

func ParseError(err error) GrpcErr {
	customErr := customErrors.GetCustomError(err)
	var validatorErr validator.ValidationErrors
	stackTrace := httpErrors.ErrorsWithStack(err)

	if err != nil {
		switch {
		case customErrors.IsDomainError(err):
			return NewDomainGrpcError(codes.Code(customErr.Status()), customErr.Error(), stackTrace)
		case customErrors.IsApplicationError(err):
			return NewApplicationGrpcError(codes.Code(customErr.Status()), customErr.Error(), stackTrace)
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
		case customErrors.IsApiError(err):
			return NewApiGrpcError(codes.Code(customErr.Status()), customErr.Error(), stackTrace)
		case customErrors.IsCustomError(err):
			return NewGrpcError(codes.Code(customErr.Status()), codes.Code(customErr.Status()).String(), customErr.Error(), stackTrace)
		case customErrors.IsUnMarshalingError(err):
			return NewInternalServerGrpcError(customErr.Error(), stackTrace)
		case customErrors.IsMarshalingError(err):
			return NewInternalServerGrpcError(customErr.Error(), stackTrace)
		case errors.Is(err, sql.ErrNoRows):
			return NewNotFoundErrorGrpcError(err.Error(), stackTrace)
		case errors.Is(err, context.DeadlineExceeded):
			return NewGrpcError(codes.DeadlineExceeded, constants.ErrRequestTimeoutTitle, err.Error(), stackTrace)
		case errors.As(err, &validatorErr):
			return NewValidationGrpcError(validatorErr.Error(), stackTrace)
		default:
			return NewInternalServerGrpcError(err.Error(), stackTrace)
		}
	}

	return nil
}
