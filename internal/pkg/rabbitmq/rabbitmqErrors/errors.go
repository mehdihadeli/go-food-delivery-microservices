package rabbitmqErrors

import (
	"emperror.dev/errors"
)

var ErrDisconnected = errors.New("disconnected from rabbitmq, trying to reconnect")
