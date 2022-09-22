package rabbitmqErrors

import (
	"emperror.dev/errors"
	"github.com/rabbitmq/amqp091-go"
)

var (
	ErrDisconnected = errors.New("disconnected from rabbitmq, trying to reconnect")
)

const (
	ConnectionError = 1
	ChannelError    = 2
)

func IsConnectionError(err *amqp091.Error) bool {
	return errorType(err.Code) == ConnectionError
}

func IsChannelError(err *amqp091.Error) bool {
	return errorType(err.Code) == ChannelError
}

func errorType(code int) int {
	switch code {
	case
		amqp091.ContentTooLarge,    // 311
		amqp091.NoConsumers,        // 313
		amqp091.AccessRefused,      // 403
		amqp091.NotFound,           // 404
		amqp091.ResourceLocked,     // 405
		amqp091.PreconditionFailed: // 406
		return ChannelError

	case
		amqp091.ConnectionForced, // 320
		amqp091.InvalidPath,      // 402
		amqp091.FrameError,       // 501
		amqp091.SyntaxError,      // 502
		amqp091.CommandInvalid,   // 503
		amqp091.ChannelError,     // 504
		amqp091.UnexpectedFrame,  // 505
		amqp091.ResourceError,    // 506
		amqp091.NotAllowed,       // 530
		amqp091.NotImplemented,   // 540
		amqp091.InternalError:    // 541
		fallthrough

	default:
		return ConnectionError
	}
}
