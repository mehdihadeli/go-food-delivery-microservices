package customErrors

import (
	"emperror.dev/errors"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/contracts"
	errorUtils "github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils/error_utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BadRequest_Err(t *testing.T) {
	rootErr2 := NewBadRequestErrorWrap(nil, fmt.Sprintf("domain event already exists in event registry"))

	rootErr := errors.New("handling bad request errorUtils")
	badErr := NewBadRequestErrorWrap(rootErr, "this is a bad request errorUtils")
	err := errors.WithMessage(badErr, "outer errorUtils wrapper")

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
	assert.Equal(t, "this is a bad request errorUtils", customError.Message())
	assert.Equal(t, "this is a bad request errorUtils: handling bad request errorUtils", customError.Error())
	assert.NotNil(t, customError.Unwrap())
	assert.NotNil(t, customError.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))           //write errorUtils messages with stacktrace
	} else {
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false))
	}

}

func Test_NotFound_Err(t *testing.T) {
	rootErr := errors.New("handling not found errorUtils")
	notFoundErr := NewNotFoundErrorWrap(rootErr, "this is a not found errorUtils")
	err := errors.WithMessage(notFoundErr, "outer errorUtils wrapper")

	assert.True(t, IsNotFoundError(err))
	assert.True(t, IsCustomError(err))

	var notFound NotFoundError
	errors.As(err, &notFound)

	assert.Equal(t, 404, notFound.Status())
	assert.Equal(t, "this is a not found errorUtils", notFound.Message())
	assert.Equal(t, "this is a not found errorUtils: handling not found errorUtils", notFound.Error())
	assert.NotNil(t, notFound.Unwrap())
	assert.NotNil(t, notFound.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))           //write errorUtils messages with stacktrace
	} else {
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false))
	}
}

func Test_Domain_Err(t *testing.T) {
	rootErr2 := NewDomainErrorWrap(nil, fmt.Sprintf("domain event already exists in event registry"))

	rootErr := errors.New("handling domain errorUtils")
	domainErr := NewDomainErrorWithCodeWrap(rootErr, 400, "this is a domain errorUtils")
	err := errors.WithMessage(domainErr, "outer errorUtils wrapper")

	assert.True(t, IsDomainError(err))
	assert.True(t, IsDomainError(rootErr2))
	assert.True(t, IsCustomError(err))
	assert.True(t, IsCustomError(rootErr2))

	var customError CustomError
	errors.As(err, &customError)

	assert.Equal(t, 400, customError.Status())
	assert.Equal(t, "this is a domain errorUtils", customError.Message())
	assert.Equal(t, "this is a domain errorUtils: handling domain errorUtils", customError.Error())
	assert.NotNil(t, customError.Unwrap())
	assert.NotNil(t, customError.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))           //write errorUtils messages with stacktrace
	} else {
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false))
	}
}

func Test_Application_Err(t *testing.T) {
	rootErr2 := NewApplicationErrorWrap(nil, fmt.Sprintf("domain event already exists in event registry"))

	rootErr := errors.New("handling application errorUtils")
	err := NewApplicationErrorWrapWithCode(rootErr, 400, "this is a application errorUtils")
	err = errors.WithMessage(err, "outer errorUtils wrapper")

	assert.True(t, IsApplicationError(err))
	assert.True(t, IsApplicationError(rootErr2))
	assert.True(t, IsCustomError(err))
	assert.True(t, IsCustomError(rootErr2))

	var appErr ApplicationError
	errors.As(err, &appErr)

	assert.Equal(t, 400, appErr.Status())
	assert.Equal(t, "this is a application errorUtils", appErr.Message())
	assert.Equal(t, "this is a application errorUtils: handling application errorUtils", appErr.Error())
	assert.NotNil(t, appErr.Unwrap())
	assert.NotNil(t, appErr.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))
	} else {
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false))
	}
}

func Test_Internal_Server_Error(t *testing.T) {
	rootErr := errors.New("handling internal server errorUtils")
	internalServerErr := NewInternalServerErrorWrap(rootErr, "this is a internal server errorUtils")
	err := errors.WithMessage(internalServerErr, "this is a internal server errorUtils")

	assert.True(t, IsInternalServerError(err))
	assert.True(t, IsCustomError(err))

	var internalErr InternalServerError
	errors.As(err, &internalErr)

	assert.Equal(t, 500, internalErr.Status())
	assert.Equal(t, "this is a internal server errorUtils", internalErr.Message())
	assert.Equal(t, "this is a internal server errorUtils: handling internal server errorUtils", internalErr.Error())
	assert.NotNil(t, internalErr.Unwrap())
	assert.NotNil(t, internalErr.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))
	} else {
		fmt.Println(errorUtils.ErrorsWithStack(err))
	}
}

func Test_Marshaling_Error(t *testing.T) {
	rootErr2 := NewMarshalingErrorWrap(nil, fmt.Sprintf("domain event already exists in event registry"))

	rootErr := errors.New("handling marshaling errorUtils")
	marshalErr := NewMarshalingErrorWrap(rootErr, "this is a marshaling errorUtils")
	err := errors.WithMessage(marshalErr, "this is a marshaling errorUtils")

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
	assert.Equal(t, "this is a marshaling errorUtils", customErr.Message())
	assert.Equal(t, "this is a marshaling errorUtils: handling marshaling errorUtils", customErr.Error())
	assert.NotNil(t, customErr.Unwrap())
	assert.NotNil(t, customErr.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))
	} else {
		fmt.Println(errorUtils.ErrorsWithStack(err))
	}
}

func Test_Validation_Error(t *testing.T) {
	rootErr2 := NewValidationErrorWrap(nil, fmt.Sprintf("domain event already exists in event registry"))

	rootErr := errors.New("handling validation errorUtils")
	validationErr := NewValidationErrorWrap(rootErr, "this is a validation errorUtils")
	err := errors.WithMessage(validationErr, "this is a validation errorUtils")

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
	assert.Equal(t, "this is a validation errorUtils", customErr.Message())
	assert.Equal(t, "this is a validation errorUtils: handling validation errorUtils", customErr.Error())
	assert.NotNil(t, customErr.Unwrap())
	assert.NotNil(t, customErr.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))
	} else {
		fmt.Println(errorUtils.ErrorsWithStack(err))
	}
}

func Test_Conflict_Error(t *testing.T) {
	rootErr2 := NewConflictErrorWrap(nil, fmt.Sprintf("domain event already exists in event registry"))

	rootErr := errors.New("handling conflict errorUtils")
	conflictErr := NewConflictErrorWrap(rootErr, "this is a conflict errorUtils")
	err := errors.WithMessage(conflictErr, "this is a conflict errorUtils")

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
	assert.Equal(t, "this is a conflict errorUtils", customErr.Message())
	assert.Equal(t, "this is a conflict errorUtils: handling conflict errorUtils", customErr.Error())
	assert.NotNil(t, customErr.Unwrap())
	assert.NotNil(t, customErr.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))
	} else {
		fmt.Println(errorUtils.ErrorsWithStack(err))
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
