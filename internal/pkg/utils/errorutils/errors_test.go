package errorUtils

import (
	"fmt"
	"testing"

	"emperror.dev/errors"
	"github.com/stretchr/testify/assert"
)

func Test_StackTraceWithErrors(t *testing.T) {
	err := errors.WithMessage(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := ErrorsWithStack(err)
	fmt.Println(res)
}

func Test_StackTrace(t *testing.T) {
	err := errors.WithMessage(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := StackTrace(err)
	fmt.Println(res)
}

func Test_RootStackTrace(t *testing.T) {
	err := errors.WithMessage(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := RootStackTrace(err)
	fmt.Println(res)
}

func Test_All_Level_Stack_Trace(t *testing.T) {
	err := errors.WrapIf(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := ErrorsWithStack(err)
	fmt.Println(res)
}

func Test_Errors_Without_Stack_Trace(t *testing.T) {
	err := errors.WrapIf(errors.New("handling bad request"), "this is a bad-request")
	err = errors.WrapIf(err, "outer error message")

	res := ErrorsWithoutStack(err, true)
	fmt.Println(res)
	assert.Contains(t, res, "outer error message\nthis is a bad-request\nhandling bad request")

	res = ErrorsWithoutStack(err, false)
	fmt.Println(res)
	assert.Contains(t, res, "outer error message: this is a bad-request: handling bad request")
}
