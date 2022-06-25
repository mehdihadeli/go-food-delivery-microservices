// Ref: https://github.com/adzeitor/mediatr

package mediatr

import (
	"context"
	"fmt"
	"reflect"
)

type Mediator struct {
	subscriptions map[reflect.Type][]reflect.Value
	registrations map[reflect.Type]reflect.Value
}

// New return new instance of mediator.
func New() Mediator {
	return Mediator{
		subscriptions: make(map[reflect.Type][]reflect.Value),
		registrations: make(map[reflect.Type]reflect.Value),
	}
}

type Handler interface {
	Handle(context.Context, interface{}) error
}

// Register registers command handler.
// Command type is detected by argument of handler.
func (m Mediator) Register(handlers ...interface{}) error {

	for _, handler := range handlers {
		handlerValue := reflect.ValueOf(handler)
		var handleMethod = handlerValue.MethodByName("Handle")

		messageType := handleMethod.Type().In(0)
		fmt.Print(messageType)
		if handleMethod.Type().NumIn() > 1 {
			if argIsContext(messageType) {
				messageType = handleMethod.Type().In(1)
			}
		}

		_, exist := m.registrations[messageType]
		if exist {
			return fmt.Errorf("handler already registered for message %T", messageType)
		}

		m.registrations[messageType] = handlerValue
	}

	return nil
}

func (m Mediator) Send(ctx context.Context, command interface{}) (interface{}, error) {
	handler, ok := m.registrations[reflect.TypeOf(command)]
	if !ok {
		return nil, fmt.Errorf("no handlers for command %T", command)
	}

	arguments := []reflect.Value{
		reflect.ValueOf(command),
	}

	var handleMethod = handler.MethodByName("Handle")

	if handleMethod.Type().NumIn() == 2 {
		arguments = append(
			[]reflect.Value{reflect.ValueOf(ctx)},
			arguments...,
		)
	}

	result := handleMethod.Call(arguments)
	switch len(result) {
	case 0:
		return nil, nil
	case 1:
		return oneReturnValuesCommand(result)
	case 2:
		return twoReturnValuesCommand(result)
	}
	return nil, nil
}

func oneReturnValuesCommand(result []reflect.Value) (interface{}, error) {
	err, isError := result[0].Interface().(error)
	if isError {
		return nil, err
	}
	return result[0].Interface(), nil
}

func twoReturnValuesCommand(result []reflect.Value) (interface{}, error) {
	var err error
	if !result[1].IsNil() {
		err = result[1].Interface().(error)
	}
	return result[0].Interface(), err
}

var contextType = reflect.TypeOf(new(context.Context)).Elem()

func argIsContext(typeOf reflect.Type) bool {
	return contextType == typeOf
}
