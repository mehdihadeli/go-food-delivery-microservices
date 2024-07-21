package contracts

import (
	"fmt"

	"emperror.dev/errors"
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

type BaseError interface {
	error
	Wrapper
	Causer
	Formatter
}
