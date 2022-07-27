package utils

import (
	"github.com/pkg/errors"
	"strings"
)

var ErrInvalidCommand = errors.New("invalid command")

func CheckType(ok bool) error {
	if !ok {
		return errors.Wrap(ErrInvalidCommand, "failed command assertion")
	}
	return nil
}

func CheckErrMessages(err error, messages ...string) bool {
	for _, message := range messages {
		if strings.Contains(strings.TrimSpace(strings.ToLower(err.Error())), strings.TrimSpace(strings.ToLower(message))) {
			return true
		}
	}
	return false
}
