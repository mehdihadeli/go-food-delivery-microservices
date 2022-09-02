package customErrors

import (
	"fmt"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors/contracts"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BadRequest_Err(t *testing.T) {
	d := errors.New("handling bad request error")
	err := NewBadRequestErrorWrap(d, "this is a bad request error")
	err = errors.WithMessage(err, "outer error wrapper")

	if IsBadRequestError(err) {
		var bad BadRequestError
		errors.As(err, &bad)

		customError := bad.GetCustomError()
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
	} else {
		assert.Fail(t, "assert to bad request failed")
	}
}

func Test_NotFound_Err(t *testing.T) {
	d := errors.New("handling not found error")
	err := NewNotFoundErrorWrap(d, "this is a not found error")
	err = errors.WithMessage(err, "outer error wrapper")

	if IsNotFoundError(err) {
		var notFound NotFoundError
		errors.As(err, &notFound)

		customError := notFound.GetCustomError()
		assert.Equal(t, 404, customError.Status())
		assert.Equal(t, "this is a not found error", customError.Message())
		assert.Equal(t, "this is a not found error: handling not found error", customError.Error())
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
	} else {
		assert.Fail(t, "assert to not found failed")
	}
}

func Test_Domain_Err(t *testing.T) {
	d := errors.New("handling domain error")
	err := NewDomainErrorWrap(d, 400, "this is a domain error")
	err = errors.WithMessage(err, "outer error wrapper")

	if IsDomainError(err) {
		var domainErr DomainError
		errors.As(err, &domainErr)

		customError := domainErr.GetCustomError()
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
	} else {
		assert.Fail(t, "assert to domain error failed")
	}
}

func Test_Application_Err(t *testing.T) {
	d := errors.New("handling application error")
	err := NewApplicationErrorWrap(d, 400, "this is a application error")
	err = errors.WithMessage(err, "outer error wrapper")

	if IsApplicationError(err) {
		var appErr ApplicationError
		errors.As(err, &appErr)

		customError := appErr.GetCustomError()
		assert.Equal(t, 400, customError.Status())
		assert.Equal(t, "this is a application error", customError.Message())
		assert.Equal(t, "this is a application error: handling application error", customError.Error())
		assert.NotNil(t, customError.Unwrap())
		assert.NotNil(t, customError.Cause())

		var stackErr contracts.StackTracer
		if ok := errors.As(err, &stackErr); ok {
			//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
			fmt.Println(httpErrors.ErrorsWithoutStack(err, false)) // Just write error messages for
			fmt.Println(httpErrors.ErrorsWithStack(err))
		} else {
			fmt.Println(httpErrors.ErrorsWithoutStack(err, false))
		}
	} else {
		assert.Fail(t, "assert to application error failed")
	}
}

func Test_Internal_Server_Error(t *testing.T) {
	d := errors.New("handling internal server error")
	err := NewInternalServerErrorWrap(d, "this is a internal server error")
	err = errors.WithMessage(err, "this is a internal server error")

	if IsInternalServerError(err) {
		var internalErr CustomError
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
	} else {
		assert.Fail(t, "assert to internal error failed")
	}
}

func Test_Marshaling_Error(t *testing.T) {
	d := errors.New("handling marshaling error")
	err := NewMarshalingErrorWrap(d, "this is a marshaling error")
	err = errors.WithMessage(err, "this is a marshaling error")

	if IsMarshalingError(err) && IsInternalServerError(err) {
		var marshalingErr MarshalingError
		errors.As(err, &marshalingErr)

		customErr := marshalingErr.GetCustomError()

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
	} else {
		assert.Fail(t, "assert to marshaling error failed")
	}
}

func myfoo(e error) error {
	//https://itnext.io/golang-error-handling-best-practice-a36f47b0b94c
	//Note: Do not repeat Wrap, it will record redundancy call stacks, we usually care about root stack trace
	return errors.WithMessage(e, "foo failed") // or errors.WrapIf()
}

func mybar(e error) error {
	return errors.WithMessage(myfoo(e), "bar failed") // or errors.WrapIf()
}
