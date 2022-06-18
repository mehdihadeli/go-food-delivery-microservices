package errors

import "github.com/pkg/errors"

var ErrInvalidCommand = errors.New("invalid command")

func CheckType(ok bool) error {
	if !ok {
		return errors.Wrap(ErrInvalidCommand, "failed command assertion")
	}
	return nil
}
