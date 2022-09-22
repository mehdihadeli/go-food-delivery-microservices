package utils

import (
	"github.com/iancoleman/strcase"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	"reflect"
)

func GetMessageName(message any) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}
	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

func GetTopicOrExchangeName(message types.IMessage) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}
	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

func GetQueueName(message types.IMessage) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}
	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

func GetRoutingKey(message types.IMessage) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}
	return strcase.ToSnake(reflect.TypeOf(message).Name())
}
