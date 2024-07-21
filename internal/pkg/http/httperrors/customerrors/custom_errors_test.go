package customErrors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/contracts"
	errorUtils "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils/errorutils"

	"emperror.dev/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Domain_Err(t *testing.T) {
	rootErr2 := NewDomainErrorWrap(
		nil,
		fmt.Sprintf("domain_events event already exists in event registry"),
	)

	// `NewPlain` doesn't add stack-trace but `New` will add stack-trace
	rootErr := errors.NewPlain("handling domain_events errorUtils")
	domainErr := NewDomainErrorWithCodeWrap(rootErr, 400, "this is a domain_events errorUtils")
	err := errors.WithMessage(domainErr, "outer errorUtils wrapper")

	assert.True(t, IsDomainError(err, 400))
	assert.True(t, IsDomainError(rootErr2, 400))
	assert.True(t, IsCustomError(err))
	assert.True(t, IsCustomError(rootErr2))

	var domainError DomainError
	errors.As(err, &domainError)

	_, isConflict := domainErr.(ConflictError)
	_, isConflict2 := domainError.(ConflictError)
	assert.False(t, isConflict)
	assert.False(t, isConflict2)

	_, isDomainError := domainErr.(DomainError)
	_, isDomainError2 := domainError.(DomainError)
	assert.True(t, isDomainError)
	assert.True(t, isDomainError2)

	assert.True(t, IsDomainError(domainErr, 400))
	assert.True(t, IsDomainError(domainError, 400))
	assert.False(t, IsDomainError(NewConflictError("conflict error"), 400))

	assert.Equal(t, 400, domainError.Status())
	assert.Equal(t, "this is a domain_events errorUtils", domainError.Message())
	assert.Equal(
		t,
		"this is a domain_events errorUtils: domain error: handling domain_events errorUtils",
		domainError.Error(),
	)
	assert.NotNil(t, domainError.Unwrap())
	assert.NotNil(t, domainError.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		// https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(
			errorUtils.ErrorsWithStack(err),
		) // write errorUtils messages with stacktrace
	} else {
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false))
	}
}

func Test_Application_Err(t *testing.T) {
	rootErr2 := NewApplicationErrorWrap(
		nil,
		fmt.Sprintf("domain_events event already exists in event registry"),
	)

	// `NewPlain` doesn't add stack-trace but `New` will add stack-trace
	rootErr := errors.NewPlain("handling application_exceptions errorUtils")
	appErr := NewApplicationErrorWrapWithCode(
		rootErr,
		400,
		"this is a application_exceptions errorUtils",
	)
	err := errors.WithMessage(appErr, "outer errorUtils wrapper")

	assert.True(t, IsApplicationError(err, 400))
	assert.True(t, IsApplicationError(rootErr2, 500))
	assert.True(t, IsCustomError(err))
	assert.True(t, IsCustomError(rootErr2))

	var applicationError ApplicationError
	errors.As(err, &applicationError)

	_, isConflict := appErr.(ConflictError)
	_, isConflict2 := applicationError.(ConflictError)
	assert.False(t, isConflict)
	assert.False(t, isConflict2)

	_, isApplicationError := appErr.(ApplicationError)
	_, isApplicationError2 := applicationError.(ApplicationError)
	assert.True(t, isApplicationError)
	assert.True(t, isApplicationError2)

	assert.True(t, IsApplicationError(appErr, 400))
	assert.True(t, IsApplicationError(applicationError, 400))
	assert.False(t, IsApplicationError(NewConflictError("conflict error"), 400))

	assert.Equal(t, 400, applicationError.Status())
	assert.Equal(t, "this is a application_exceptions errorUtils", applicationError.Message())
	assert.Equal(
		t,
		"this is a application_exceptions errorUtils: application error: handling application_exceptions errorUtils",
		applicationError.Error(),
	)
	assert.NotNil(t, applicationError.Unwrap())
	assert.NotNil(t, applicationError.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		// https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))
	} else {
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false))
	}
}

func Test_Api_Err(t *testing.T) {
	rootErr2 := NewApiErrorWrap(
		nil,
		http.StatusBadRequest,
		fmt.Sprintf("domain_events event already exists in event registry"),
	)

	// `NewPlain` doesn't add stack-trace but `New` will add stack-trace
	rootErr := errors.NewPlain("handling api_exceptions errorUtils")
	appErr := NewApiErrorWrap(
		rootErr,
		400,
		"this is a api_exceptions errorUtils",
	)
	err := errors.WithMessage(appErr, "outer errorUtils wrapper")

	assert.True(t, IsApiError(err, 400))
	assert.True(t, IsApiError(rootErr2, 500))
	assert.True(t, IsCustomError(err))
	assert.True(t, IsCustomError(rootErr2))

	var apiError ApiError
	errors.As(err, &apiError)

	_, isConflict := appErr.(ConflictError)
	_, isConflict2 := apiError.(ConflictError)
	assert.False(t, isConflict)
	assert.False(t, isConflict2)

	_, isApiError := appErr.(ApiError)
	_, isApiError2 := apiError.(ApiError)
	assert.True(t, isApiError)
	assert.True(t, isApiError2)

	assert.True(t, IsApiError(appErr, 400))
	assert.True(t, IsApiError(apiError, 400))
	assert.False(t, IsApiError(NewConflictError("conflict error"), 400))

	assert.Equal(t, 400, apiError.Status())
	assert.Equal(t, "this is a api_exceptions errorUtils", apiError.Message())
	assert.Equal(
		t,
		"this is a api_exceptions errorUtils: api error: handling api_exceptions errorUtils",
		apiError.Error(),
	)
	assert.NotNil(t, apiError.Unwrap())
	assert.NotNil(t, apiError.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		// https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))
	} else {
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false))
	}
}

func Test_BadRequest_Err(t *testing.T) {
	rootErr2 := NewBadRequestErrorWrap(
		nil,
		fmt.Sprintf("domain_events event already exists in event registry"),
	)

	// `NewPlain` doesn't add stack-trace but `New` will add stack-trace
	rootErr := errors.NewPlain("handling bad request errorUtils")
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

	var badRequestError BadRequestError
	errors.As(err, &badRequestError)

	_, isConflict := badErr.(ConflictError)
	_, isConflict2 := badRequestError.(ConflictError)
	assert.False(t, isConflict)
	assert.False(t, isConflict2)

	_, isBadRequest := badErr.(BadRequestError)
	_, isBadRequest2 := badRequestError.(BadRequestError)
	assert.True(t, isBadRequest)
	assert.True(t, isBadRequest2)

	assert.True(t, IsBadRequestError(badErr))
	assert.True(t, IsBadRequestError(badRequestError))
	assert.False(t, IsBadRequestError(NewConflictError("conflict error")))

	assert.Equal(t, 400, badRequestError.Status())
	assert.Equal(t, "this is a bad request errorUtils", badRequestError.Message())
	assert.Equal(
		t,
		"this is a bad request errorUtils: bad request error: handling bad request errorUtils",
		badRequestError.Error(),
	)
	assert.NotNil(t, badRequestError.Unwrap())
	assert.NotNil(t, badRequestError.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		// https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(
			errorUtils.ErrorsWithStack(err),
		) // write errorUtils messages with stacktrace
	} else {
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false))
	}
}

func Test_NotFound_Err(t *testing.T) {
	// `NewPlain` doesn't add stack-trace but `New` will add stack-trace
	rootErr := errors.NewPlain("handling not found errorUtils")
	notFoundErr := NewNotFoundErrorWrap(rootErr, "this is a not found errorUtils")
	err := errors.WithMessage(notFoundErr, "outer errorUtils wrapper")

	assert.True(t, IsNotFoundError(err))
	assert.True(t, IsCustomError(err))

	var notFound NotFoundError
	errors.As(err, &notFound)

	_, isConflict := notFoundErr.(ConflictError)
	_, isConflict2 := notFound.(ConflictError)
	assert.False(t, isConflict)
	assert.False(t, isConflict2)

	_, isNotFound := notFoundErr.(NotFoundError)
	_, isNotFound2 := notFound.(NotFoundError)
	assert.True(t, isNotFound)
	assert.True(t, isNotFound2)

	assert.True(t, IsNotFoundError(notFoundErr))
	assert.True(t, IsNotFoundError(notFound))
	assert.False(t, IsNotFoundError(NewConflictError("conflict error")))

	assert.Equal(t, http.StatusNotFound, notFound.Status())
	assert.Equal(t, "this is a not found errorUtils", notFound.Message())
	assert.Equal(
		t,
		"this is a not found errorUtils: not found error: handling not found errorUtils",
		notFound.Error(),
	)
	assert.NotNil(t, notFound.Unwrap())
	assert.NotNil(t, notFound.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		// https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(
			errorUtils.ErrorsWithStack(err),
		) // write errorUtils messages with stacktrace
	} else {
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false))
	}
}

func Test_Internal_Server_Error(t *testing.T) {
	// `NewPlain` doesn't add stack-trace but `New` will add stack-trace
	rootErr := errors.NewPlain("handling internal server errorUtils")
	internalServerErr := NewInternalServerErrorWrap(rootErr, "this is a internal server errorUtils")
	err := errors.WithMessage(internalServerErr, "outer errorUtils wrapper")

	assert.True(t, IsInternalServerError(err))
	assert.True(t, IsCustomError(err))

	var internalErr InternalServerError
	errors.As(err, &internalErr)

	assert.True(t, IsInternalServerError(internalErr))
	assert.False(t, IsInternalServerError(NewConflictError("conflict error")))

	assert.Equal(t, http.StatusInternalServerError, internalErr.Status())
	assert.Equal(t, "this is a internal server errorUtils", internalErr.Message())
	assert.Equal(
		t,
		"this is a internal server errorUtils: internal server error: handling internal server errorUtils",
		internalErr.Error(),
	)
	assert.NotNil(t, internalErr.Unwrap())
	assert.NotNil(t, internalErr.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		// https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))
	} else {
		fmt.Println(errorUtils.ErrorsWithStack(err))
	}
}

func Test_Forbidden_Error(t *testing.T) {
	// `NewPlain` doesn't add stack-trace but `New` will add stack-trace
	rootErr := errors.NewPlain("handling forbidden errorUtils")
	forbiddenError := NewForbiddenErrorWrap(rootErr, "this is a forbidden errorUtils")
	err := errors.WithMessage(forbiddenError, "outer errorUtils wrapper")

	assert.True(t, IsForbiddenError(err))
	assert.True(t, IsCustomError(err))

	var forbiddenErr ForbiddenError
	errors.As(err, &forbiddenErr)

	assert.True(t, IsForbiddenError(forbiddenErr))
	assert.False(t, IsForbiddenError(NewConflictError("conflict error")))

	assert.Equal(t, http.StatusForbidden, forbiddenErr.Status())
	assert.Equal(t, "this is a forbidden errorUtils", forbiddenErr.Message())
	assert.Equal(
		t,
		"this is a forbidden errorUtils: forbidden error: handling forbidden errorUtils",
		forbiddenErr.Error(),
	)
	assert.NotNil(t, forbiddenErr.Unwrap())
	assert.NotNil(t, forbiddenErr.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		// https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))
	} else {
		fmt.Println(errorUtils.ErrorsWithStack(err))
	}
}

func Test_Marshaling_Error(t *testing.T) {
	rootErr2 := NewMarshalingErrorWrap(
		nil,
		fmt.Sprintf("domain_events event already exists in event registry"),
	)

	rootErr := errors.NewPlain("handling marshaling errorUtils")
	marshalErr := NewMarshalingErrorWrap(rootErr, "this is a marshaling errorUtils")
	err := errors.WithMessage(marshalErr, "outer errorUtils wrapper")

	assert.True(t, IsInternalServerError(err))
	assert.True(t, IsInternalServerError(rootErr2))
	assert.True(t, IsCustomError(err))
	assert.True(t, IsCustomError(rootErr2))

	assert.True(t, IsMarshalingError(err))
	assert.True(t, IsMarshalingError(rootErr2))

	var marshallingErr MarshalingError
	errors.As(err, &marshallingErr)

	assert.True(t, IsMarshalingError(marshallingErr))
	assert.False(t, IsMarshalingError(NewConflictError("conflict error")))

	assert.Equal(t, 500, marshallingErr.Status())
	assert.Equal(t, "this is a marshaling errorUtils", marshallingErr.Message())
	assert.Equal(
		t,
		"this is a marshaling errorUtils: marshaling error: handling marshaling errorUtils",
		marshallingErr.Error(),
	)
	assert.NotNil(t, marshallingErr.Unwrap())
	assert.NotNil(t, marshallingErr.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		// https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))
	} else {
		fmt.Println(errorUtils.ErrorsWithStack(err))
	}
}

func Test_Validation_Error(t *testing.T) {
	rootErr2 := NewValidationErrorWrap(
		nil,
		fmt.Sprintf("domain_events event already exists in event registry"),
	)

	rootErr := errors.New("handling validation errorUtils")
	validationErr := NewValidationErrorWrap(rootErr, "this is a validation errorUtils")
	err := errors.WithMessage(validationErr, "this is a top error message")

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

	assert.Equal(t, http.StatusBadRequest, customErr.Status())
	assert.Equal(t, "this is a validation errorUtils", customErr.Message())
	assert.Equal(
		t,
		"this is a validation errorUtils: validation error: handling validation errorUtils",
		customErr.Error(),
	)
	assert.NotNil(t, customErr.Unwrap())
	assert.NotNil(t, customErr.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		// https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))
	} else {
		fmt.Println(errorUtils.ErrorsWithStack(err))
	}
}

func Test_Conflict_Error(t *testing.T) {
	rootErr2 := NewConflictErrorWrap(
		nil,
		fmt.Sprintf("domain_events event already exists in event registry"),
	)

	// `NewPlain` doesn't add stack-trace but `New` will add stack-trace
	rootErr := errors.NewPlain("handling conflict errorUtils")
	conflictErr := NewConflictErrorWrap(rootErr, "this is a conflict errorUtils")
	err := errors.WithMessage(conflictErr, "this is a top error message")

	assert.True(t, IsCustomError(err))
	assert.True(t, IsConflictError(err))
	assert.True(t, IsCustomError(rootErr2))
	assert.True(t, IsConflictError(rootErr2))

	var conflictError ConflictError
	errors.As(err, &conflictError)

	assert.Equal(t, 409, conflictError.Status())
	assert.Equal(t, "this is a conflict errorUtils", conflictError.Message())
	assert.Equal(
		t,
		"this is a conflict errorUtils: conflict error: handling conflict errorUtils",
		conflictError.Error(),
	)
	assert.NotNil(t, conflictError.Unwrap())
	assert.NotNil(t, conflictError.Cause())

	var stackErr contracts.StackTracer
	if ok := errors.As(err, &stackErr); ok {
		// https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
		fmt.Println(errorUtils.ErrorsWithoutStack(err, false)) // Just write errorUtils messages for
		fmt.Println(errorUtils.ErrorsWithStack(err))
	} else {
		fmt.Println(errorUtils.ErrorsWithStack(err))
	}
}

func myfoo(e error) error {
	// https://itnext.io/golang-error-handling-best-practice-a36f47b0b94c
	// Note: Do not repeat Wrap, it will record redundancy call stacks, we usually care about root stack trace
	return errors.WithMessage(e, "foo failed") // or grpc_errors.WrapIf()
}

func mybar(e error) error {
	return errors.WithMessage(myfoo(e), "bar failed") // or grpc_errors.WrapIf()
}
