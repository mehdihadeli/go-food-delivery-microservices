package problemDetails

import (
	"context"
	"database/sql"
	"net/http"
	"reflect"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/constants"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"
	errorUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/error_utils"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
)

type ProblemDetailParser struct {
	internalErrors map[reflect.Type]func(err error) ProblemDetailErr
}

func NewProblemDetailParser(builder func(builder *OptionBuilder)) *ProblemDetailParser {
	optionBuilder := NewOptionBuilder()
	builder(optionBuilder)
	items := optionBuilder.Build()
	return &ProblemDetailParser{internalErrors: items}
}

func (p *ProblemDetailParser) ResolveError(err error) ProblemDetailErr {
	errType := typeMapper.GetReflectType(err)
	problem := p.internalErrors[errType]
	if problem != nil {
		return problem(err)
	}
	return nil
}

func ParseError(err error) ProblemDetailErr {
	stackTrace := errorUtils.ErrorsWithStack(err)
	customErr := customErrors.GetCustomError(err)
	var validatorErr validator.ValidationErrors

	if err != nil {
		switch {
		case customErrors.IsDomainError(err, customErr.Status()):
			return NewDomainProblemDetail(customErr.Status(), customErr.Error(), stackTrace)
		case customErrors.IsApplicationError(err, customErr.Status()):
			return NewApplicationProblemDetail(customErr.Status(), customErr.Error(), stackTrace)
		case customErrors.IsApiError(err, customErr.Status()):
			return NewApiProblemDetail(customErr.Status(), customErr.Error(), stackTrace)
		case customErrors.IsBadRequestError(err):
			return NewBadRequestProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsNotFoundError(err):
			return NewNotFoundErrorProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsValidationError(err):
			return NewValidationProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsUnAuthorizedError(err):
			return NewUnAuthorizedErrorProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsForbiddenError(err):
			return NewForbiddenProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsConflictError(err):
			return NewConflictProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsInternalServerError(err):
			return NewInternalServerProblemDetail(customErr.Error(), stackTrace)
		case customErrors.IsCustomError(err):
			return NewProblemDetailFromCodeAndDetail(customErr.Status(), customErr.Error(), stackTrace)
		case customErrors.IsUnMarshalingError(err):
			return NewInternalServerProblemDetail(err.Error(), stackTrace)
		case customErrors.IsMarshalingError(err):
			return NewInternalServerProblemDetail(err.Error(), stackTrace)
		case errors.Is(err, sql.ErrNoRows):
			return NewNotFoundErrorProblemDetail(err.Error(), stackTrace)
		case errors.Is(err, context.DeadlineExceeded):
			return NewProblemDetail(
				http.StatusRequestTimeout,
				constants.ErrRequestTimeoutTitle,
				err.Error(),
				stackTrace,
			)
		case errors.As(err, &validatorErr):
			return NewValidationProblemDetail(validatorErr.Error(), stackTrace)
		default:
			return NewInternalServerProblemDetail(err.Error(), stackTrace)
		}
	}

	return nil
}
