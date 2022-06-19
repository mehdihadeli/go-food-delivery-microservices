package es

import "github.com/pkg/errors"

var (
	ErrAlreadyExists       = errors.New("Already exists")
	ErrAggregateNotFound   = errors.New("aggregate not found")
	ErrInvalidEventType    = errors.New("invalid event type")
	ErrInvalidCommandType  = errors.New("invalid command type")
	ErrInvalidAggregate    = errors.New("invalid aggregate")
	ErrInvalidAggregateID  = errors.New("invalid aggregate id")
	ErrInvalidEventVersion = errors.New("invalid event version")
)
