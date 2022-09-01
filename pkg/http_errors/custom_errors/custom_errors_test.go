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
	s := errors.New("handling bad request")
	err := NewBadRequestErrorWrap(s, "this is a bad-request").WithStack()

	if IsBadRequestError(err) {
		var bad BadRequestError
		errors.As(err, &bad)

		assert.Equal(t, 400, bad.Status())
		assert.Equal(t, "this is a bad-request", bad.Message())
		assert.Equal(t, "this is a bad-request: handling bad request", bad.Error())
		assert.NotNil(t, bad.Unwrap())
		assert.NotNil(t, bad.Cause())

		var stackErr contracts.StackTracer
		if ok := errors.As(err, &stackErr); ok {
			//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
			fmt.Println(httpErrors.ErrorsWithoutTrace(err, false)) // Just write error messages for
			fmt.Println(httpErrors.ErrorsWithStack(err))           //write error messages with stacktrace
		} else {
			fmt.Println(httpErrors.ErrorsWithoutTrace(err, false))
		}
	} else {
		assert.Fail(t, "assert to bad request failed")
	}
}

func Test_NotFound_Err(t *testing.T) {
	s := errors.New("handling bad request")
	err := NewNotFoundErrorWrap(s, "this is a not found").WithStack()
	newError := mybar(err)

	g := errors.WithStack(NewCustomError(s, 400, ""))
	h, ok := g.(CustomError)
	fmt.Println(ok)
	fmt.Println(h)

	if IsNotFoundError(newError) {
		var notFound NotFoundError
		errors.As(newError, &notFound)

		assert.Equal(t, 404, notFound.Status())
		assert.Equal(t, "this is a not found", notFound.Message())
		assert.Equal(t, "this is a not found: handling bad request", notFound.Error())
		assert.NotNil(t, notFound.Unwrap())
		assert.NotNil(t, notFound.Cause())

		var stackErr contracts.StackTracer
		if ok := errors.As(newError, &stackErr); ok {
			//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
			//fmt.Printf("%v\n\n", err.Error())
			fmt.Println(httpErrors.ErrorsWithoutTrace(newError, false)) // Just write error messages for
			fmt.Println(httpErrors.ErrorsWithStack(newError))           //write error messages with stacktrace
		} else {
			fmt.Println(httpErrors.ErrorsWithoutTrace(newError, false))
		}
	} else {
		assert.Fail(t, "assert to not found failed")
	}
}

func Test_Domain_Err(t *testing.T) {
	s := errors.New("handling domain error")
	err := NewDomainErrorWrap(s, 400, "this is a domain error").WithStack()
	newError := mybar(err)

	if IsDomainError(newError) {
		var domainErr DomainError
		errors.As(newError, &domainErr)

		assert.Equal(t, 400, domainErr.Status())
		assert.Equal(t, "this is a domain error", domainErr.Message())
		assert.Equal(t, "this is a domain error: handling domain error", domainErr.Error())
		assert.NotNil(t, domainErr.Unwrap())
		assert.NotNil(t, domainErr.Cause())

		var stackErr contracts.StackTracer
		if ok := errors.As(newError, &stackErr); ok {
			//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
			fmt.Println(httpErrors.ErrorsWithoutTrace(newError, false)) // Just write error messages for
			fmt.Println(httpErrors.ErrorsWithStack(newError))           //write error messages with stacktrace
		} else {
			fmt.Println(httpErrors.ErrorsWithoutTrace(newError, false))
		}
	} else {
		assert.Fail(t, "assert to domain error failed")
	}
}

func Test_Application_Err(t *testing.T) {
	s := errors.New("handling application error")
	err := NewApplicationErrorWrap(s, 400, "this is a application error").WithStack()
	newError := mybar(err)

	if IsApplicationError(newError) {
		var appErr ApplicationError
		errors.As(newError, &appErr)

		assert.Equal(t, 400, appErr.Status())
		assert.Equal(t, "this is a application error", appErr.Message())
		assert.Equal(t, "this is a application error: handling application error", appErr.Error())
		assert.NotNil(t, appErr.Unwrap())
		assert.NotNil(t, appErr.Cause())

		var stackErr contracts.StackTracer
		if ok := errors.As(newError, &stackErr); ok {
			//https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
			fmt.Println(httpErrors.ErrorsWithoutTrace(newError, false)) // Just write error messages for
			fmt.Println(httpErrors.ErrorsWithStack(newError))
		} else {
			fmt.Println(httpErrors.ErrorsWithoutTrace(newError, false))
		}
	} else {
		assert.Fail(t, "assert to application error failed")
	}
}

func Test_Internal_Server_error(t *testing.T) {
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
			fmt.Println(httpErrors.ErrorsWithoutTrace(err, false)) // Just write error messages for
			fmt.Println(httpErrors.ErrorsWithStack(err))
		} else {
			fmt.Println(httpErrors.ErrorsWithStack(err))
		}
	} else {
		assert.Fail(t, "assert to internal error failed")
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
