package customErrors

import (
	"emperror.dev/errors"
	"fmt"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/contracts"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BadRequest_Err(t *testing.T) {
	rootErr2 := NewBadRequestErrorWrap(nil, fmt.Sprintf("domain event already exists in event registry"))

	rootErr := errors.New("handling bad request error")
	badErr := NewBadRequestErrorWrap(rootErr, "this is a bad request error")
	err := errors.WithMessage(badErr, "outer error wrapper")

	assert.True(t, IsBadRequestError(err))
	assert.True(t, IsCustomError(err))
	assert.True(t, IsCustomError(rootErr2))
	assert.True(t, IsCustomError(rootErr2))

	var customError CustomError
	var customError2 CustomError
	errors.As(err, &customError)
	errors.As(err, &customError2)

	assert.NotNil(t, customError2)

	assert.Equal(t, 400, customError.Status())
	assert.Equal(t, "this is a bad request error", customError.Message())
	assert.Equal(t, "this is a bad request error: handling bad request error", customError.Error())
	assert.NotNil(t, customError.Unwrap())
	assert.NotNil(t, customError.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(httpErrors.ErrorsWithoutStack(err, false)) // Just write error messages for
		fmt.Println(httpErrors.ErrorsWithStack(err))           //write error messages with stacktrace
	} else {
		fmt.Println(httpErrors.ErrorsWithoutStack(err, false))
	}

}

func Test_NotFound_Err(t *testing.T) {
	rootErr := errors.New("handling not found error")
	notFoundErr := NewNotFoundErrorWrap(rootErr, "this is a not found error")
	err := errors.WithMessage(notFoundErr, "outer error wrapper")

	assert.True(t, IsNotFoundError(err))
	assert.True(t, IsCustomError(err))

	var notFound NotFoundError
	errors.As(err, &notFound)

	assert.Equal(t, 404, notFound.Status())
	assert.Equal(t, "this is a not found error", notFound.Message())
	assert.Equal(t, "this is a not found error: handling not found error", notFound.Error())
	assert.NotNil(t, notFound.Unwrap())
	assert.NotNil(t, notFound.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(httpErrors.ErrorsWithoutStack(err, false)) // Just write error messages for
		fmt.Println(httpErrors.ErrorsWithStack(err))           //write error messages with stacktrace
	} else {
		fmt.Println(httpErrors.ErrorsWithoutStack(err, false))
	}
}

func Test_Domain_Err(t *testing.T) {
	rootErr2 := NewDomainErrorWrap(nil, fmt.Sprintf("domain event already exists in event registry"))

	rootErr := errors.New("handling domain error")
	domainErr := NewDomainErrorWithCodeWrap(rootErr, 400, "this is a domain error")
	err := errors.WithMessage(domainErr, "outer error wrapper")

	assert.True(t, IsDomainError(err))
	assert.True(t, IsDomainError(rootErr2))
	assert.True(t, IsCustomError(err))
	assert.True(t, IsCustomError(rootErr2))

	var customError CustomError
	errors.As(err, &customError)

	assert.Equal(t, 400, customError.Status())
	assert.Equal(t, "this is a domain error", customError.Message())
	assert.Equal(t, "this is a domain error: handling domain error", customError.Error())
	assert.NotNil(t, customError.Unwrap())
	assert.NotNil(t, customError.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(httpErrors.ErrorsWithoutStack(err, false)) // Just write error messages for
		fmt.Println(httpErrors.ErrorsWithStack(err))           //write error messages with stacktrace
	} else {
		fmt.Println(httpErrors.ErrorsWithoutStack(err, false))
	}
}

func Test_Application_Err(t *testing.T) {
	rootErr2 := NewApplicationErrorWrap(nil, fmt.Sprintf("domain event already exists in event registry"))
	
	rootErr := errors.New("handling application error")
	err := NewApplicationErrorWrapWithCode(rootErr, 400, "this is a application error")
	err = errors.WithMessage(err, "outer error wrapper")

	assert.True(t, IsApplicationError(err))
	assert.True(t, IsApplicationError(rootErr2))
	assert.True(t, IsCustomError(err))
	assert.True(t, IsCustomError(rootErr2))

	var appErr ApplicationError
	errors.As(err, &appErr)

	assert.Equal(t, 400, appErr.Status())
	assert.Equal(t, "this is a application error", appErr.Message())
	assert.Equal(t, "this is a application error: handling application error", appErr.Error())
	assert.NotNil(t, appErr.Unwrap())
	assert.NotNil(t, appErr.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(httpErrors.ErrorsWithoutStack(err, false)) // Just write error messages for
		fmt.Println(httpErrors.ErrorsWithStack(err))
	} else {
		fmt.Println(httpErrors.ErrorsWithoutStack(err, false))
	}
}

func Test_Internal_Server_Error(t *testing.T) {
	rootErr := errors.New("handling internal server error")
	internalServerErr := NewInternalServerErrorWrap(rootErr, "this is a internal server error")
	err := errors.WithMessage(internalServerErr, "this is a internal server error")

	assert.True(t, IsInternalServerError(err))
	assert.True(t, IsCustomError(err))

	var internalErr InternalServerError
	errors.As(err, &internalErr)

	assert.Equal(t, 500, internalErr.Status())
	assert.Equal(t, "this is a internal server error", internalErr.Message())
	assert.Equal(t, "this is a internal server error: handling internal server error", internalErr.Error())
	assert.NotNil(t, internalErr.Unwrap())
	assert.NotNil(t, internalErr.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(httpErrors.ErrorsWithoutStack(err, false)) // Just write error messages for
		fmt.Println(httpErrors.ErrorsWithStack(err))
	} else {
		fmt.Println(httpErrors.ErrorsWithStack(err))
	}
}

func Test_Marshaling_Error(t *testing.T) {
	rootErr2 := NewMarshalingErrorWrap(nil, fmt.Sprintf("domain event already exists in event registry"))

	rootErr := errors.New("handling marshaling error")
	marshalErr := NewMarshalingErrorWrap(rootErr, "this is a marshaling error")
	err := errors.WithMessage(marshalErr, "this is a marshaling error")

	assert.True(t, IsInternalServerError(err))
	assert.True(t, IsInternalServerError(rootErr2))
	assert.True(t, IsCustomError(err))
	assert.True(t, IsCustomError(rootErr2))

	assert.True(t, IsMarshalingError(err))
	assert.True(t, IsMarshalingError(rootErr2))

	var customErr CustomError
	var customErr2 CustomError

	errors.As(err, &customErr)
	errors.As(rootErr2, &customErr2)

	assert.NotNil(t, customErr)
	assert.NotNil(t, customErr2)

	assert.Equal(t, 500, customErr.Status())
	assert.Equal(t, "this is a marshaling error", customErr.Message())
	assert.Equal(t, "this is a marshaling error: handling marshaling error", customErr.Error())
	assert.NotNil(t, customErr.Unwrap())
	assert.NotNil(t, customErr.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(httpErrors.ErrorsWithoutStack(err, false)) // Just write error messages for
		fmt.Println(httpErrors.ErrorsWithStack(err))
	} else {
		fmt.Println(httpErrors.ErrorsWithStack(err))
	}
}

func Test_Validation_Error(t *testing.T) {
	rootErr2 := NewValidationErrorWrap(nil, fmt.Sprintf("domain event already exists in event registry"))

	rootErr := errors.New("handling validation error")
	validationErr := NewValidationErrorWrap(rootErr, "this is a validation error")
	err := errors.WithMessage(validationErr, "this is a validation error")

	assert.True(t, IsBadRequestError(err))
	assert.True(t, IsBadRequestError(rootErr2))
	assert.True(t, IsCustomError(err))
	assert.True(t, IsCustomError(rootErr2))

	assert.True(t, IsValidationError(err))
	assert.True(t, IsValidationError(rootErr2))

	var customErr CustomError
	var customErr2 CustomError

	errors.As(err, &customErr)
	errors.As(rootErr2, &customErr2)

	assert.NotNil(t, customErr)
	assert.NotNil(t, customErr2)

	assert.Equal(t, 400, customErr.Status())
	assert.Equal(t, "this is a validation error", customErr.Message())
	assert.Equal(t, "this is a validation error: handling validation error", customErr.Error())
	assert.NotNil(t, customErr.Unwrap())
	assert.NotNil(t, customErr.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(httpErrors.ErrorsWithoutStack(err, false)) // Just write error messages for
		fmt.Println(httpErrors.ErrorsWithStack(err))
	} else {
		fmt.Println(httpErrors.ErrorsWithStack(err))
	}
}

func Test_Conflict_Error(t *testing.T) {
	rootErr2 := NewConflictErrorWrap(nil, fmt.Sprintf("domain event already exists in event registry"))

	rootErr := errors.New("handling conflict error")
	conflictErr := NewConflictErrorWrap(rootErr, "this is a conflict error")
	err := errors.WithMessage(conflictErr, "this is a conflict error")

	assert.True(t, IsCustomError(err))
	assert.True(t, IsConflictError(err))
	assert.True(t, IsCustomError(rootErr2))
	assert.True(t, IsConflictError(rootErr2))

	var customErr CustomError
	var customErr2 CustomError
	errors.As(err, &customErr)
	errors.As(rootErr2, &customErr2)

	assert.NotNil(t, customErr2)

	assert.Equal(t, 409, customErr.Status())
	assert.Equal(t, "this is a conflict error", customErr.Message())
	assert.Equal(t, "this is a conflict error: handling conflict error", customErr.Error())
	assert.NotNil(t, customErr.Unwrap())
	assert.NotNil(t, customErr.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(httpErrors.ErrorsWithoutStack(err, false)) // Just write error messages for
		fmt.Println(httpErrors.ErrorsWithStack(err))
	} else {
		fmt.Println(httpErrors.ErrorsWithStack(err))
	}
}

func myfoo(e error) error {
	//https://itnext.io/golang-error-handling-best-practice-a36f47b0b94c
	//Note: Do not repeat Wrap, it will record redundancy call stacks, we usually care about root stack trace
	return errors.WithMessage(e, "foo failed") // or grpc_errors.WrapIf()
}

func mybar(e error) error {
	return errors.WithMessage(myfoo(e), "bar failed") // or grpc_errors.WrapIf()
}
