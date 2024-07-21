package cqrs

import (
	"context"

	"github.com/mehdihadeli/go-mediatr"
)

type CommandHandler[TCommand Command, TResponse any] interface {
	Handle(ctx context.Context, command TCommand) (TResponse, error)
}

type CommandHandlerVoid[TCommand Command] interface {
	Handle(ctx context.Context, command TCommand) (*mediatr.Unit, error)
}
