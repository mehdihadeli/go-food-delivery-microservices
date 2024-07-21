package customErrors

import (
	"fmt"
	"io"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/contracts"

	"emperror.dev/errors"
)

// https://klotzandrew.com/blog/error-handling-in-golang
// https://banzaicloud.com/blog/error-handling-go/
// https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully
// https://dave.cheney.net/2016/06/12/stack-traces-and-the-errors-package
// https://github.com/go-stack/stack
// https://github.com/juju/errors
// https://github.com/emperror/errors
// https://github.com/pkg/errors/issues/75
type customError struct {
	statusCode int
	message    string
	error
}

type CustomError interface {
	error
	contracts.Wrapper
	contracts.Causer
	contracts.Formatter
	isCustomError()
	Status() int
	Message() string
}

func NewCustomError(err error, code int, message string) CustomError {
	m := &customError{
		statusCode: code,
		error:      err,
		message:    message,
	}

	return m
}

func (e *customError) isCustomError() {
}

func (e *customError) Error() string {
	if e.error != nil {
		return e.error.Error()
	}

	return e.message
}

func (e *customError) Message() string {
	return e.message
}

func (e *customError) Status() int {
	return e.statusCode
}

func (e *customError) Cause() error {
	return e.error
}

func (e *customError) Unwrap() error {
	return e.error
}

func (e *customError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			//%s	error messages separated by a colon and a space (": ")
			//%q	double-quoted error messages separated by a colon and a space (": ")
			//%v	one error message per line
			//%+v	one error message per line and stack trace (if any)

			// if we have a call-stacked error, +v shows callstack for this error
			fmt.Fprintf(s, "%+v", e.Cause())
			// io.WriteString(s, e.message)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, e.Error())
	}
}

func GetCustomError(err error) CustomError {
	if IsCustomError(err) {
		var internalErr CustomError
		errors.As(err, &internalErr)

		return internalErr
	}

	return nil
}

func IsCustomError(err error) bool {
	var customErr CustomError

	_, ok := err.(CustomError)
	if ok {
		return true
	}

	// us, ok := errors.Cause(err).(ConflictError)
	if errors.As(err, &customErr) {
		return true
	}

	return false
}
