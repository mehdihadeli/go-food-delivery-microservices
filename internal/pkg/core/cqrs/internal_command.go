package cqrs

type internalCommand struct {
	Command
}

type InternalCommand interface {
	Command
	isInternalCommand()
}

func NewInternalCommandByT[T any]() InternalCommand {
	return &internalCommand{Command: NewCommandByT[T]()}
}

func (c *internalCommand) isInternalCommand() {
}

func IsInternalCommand(obj interface{}) bool {
	if _, ok := obj.(InternalCommand); ok {
		return true
	}

	return false
}
