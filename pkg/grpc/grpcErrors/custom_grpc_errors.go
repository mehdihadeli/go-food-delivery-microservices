package grpcErrors

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/constants"
	"google.golang.org/grpc/codes"
	"time"
)

func NewValidationGrpcError(detail string, stackTrace string) GrpcErr {
	validationError :=
		&grpcErr{
			Title:      constants.ErrBadRequestTitle,
			Detail:     detail,
			Status:     codes.InvalidArgument,
			Timestamp:  time.Now(),
			StackTrace: stackTrace,
		}

	return validationError
}

func NewConflictGrpcError(detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrConflictTitle,
		Detail:     detail,
		Status:     codes.AlreadyExists,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewBadRequestGrpcError(detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrBadRequestTitle,
		Detail:     detail,
		Status:     codes.InvalidArgument,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewNotFoundErrorGrpcError(detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrNotFoundTitle,
		Detail:     detail,
		Status:     codes.NotFound,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewUnAuthorizedErrorGrpcError(detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrUnauthorizedTitle,
		Detail:     detail,
		Status:     codes.Unauthenticated,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewForbiddenGrpcError(detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrForbiddenTitle,
		Detail:     detail,
		Status:     codes.PermissionDenied,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewInternalServerGrpcError(detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrInternalServerErrorTitle,
		Detail:     detail,
		Status:     codes.Internal,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewDomainGrpcError(status codes.Code, detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrDomainTitle,
		Detail:     detail,
		Status:     status,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewApplicationGrpcError(status codes.Code, detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrApplicationTitle,
		Detail:     detail,
		Status:     status,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}

func NewApiGrpcError(status codes.Code, detail string, stackTrace string) GrpcErr {
	return &grpcErr{
		Title:      constants.ErrApiTitle,
		Detail:     detail,
		Status:     status,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
	}
}
