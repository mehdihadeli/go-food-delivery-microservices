package cqrs

type command struct {
	TypeInfo
}

type Command interface {
	isCommand()

	TypeInfo
}

func NewCommandByT[T any]() Command {
	c := &command{TypeInfo: NewTypeInfoT[T]()}

	return c
}

func (c *command) isCommand() {
}

func IsCommand(obj interface{}) bool {
	if _, ok := obj.(Command); ok {
		return true
	}

	return false
}
