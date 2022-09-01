package contracts

import (
	"fmt"
	"github.com/pkg/errors"
)

type Causer interface {
	Cause() error
}

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type Wrapper interface {
	Unwrap() error
}

type Formatter interface {
	Format(f fmt.State, verb rune)
}

type StackError interface {
	WithStack() error
}

type WithStack interface {
	error
	StackTracer
	Wrapper
	Causer
	Formatter
}
