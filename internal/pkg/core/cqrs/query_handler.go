package cqrs

import (
	"context"
)

type QueryHandler[TQuery Query, TResponse any] interface {
	Handle(ctx context.Context, query TQuery) (TResponse, error)
}
