package utils

import (
	"reflect"

	"github.com/ahmetb/go-linq/v3"
	"github.com/iancoleman/strcase"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"
)

func GetMessageName(message interface{}) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}
	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

func GetMessageNameFromType(message reflect.Type) string {
	if message.Kind() == reflect.Pointer {
		return strcase.ToSnake(message.Elem().Name())
	}
	return strcase.ToSnake(message.Name())
}

func GetMessageBaseReflectTypeFromType(message reflect.Type) reflect.Type {
	return typeMapper.GetBaseReflectType(message)
}

func GetMessageBaseReflectType(message interface{}) reflect.Type {
	return typeMapper.GetBaseReflectType(message)
}

func GetTopicOrExchangeName(message interface{}) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}
	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

func GetTopicOrExchangeNameFromType(message reflect.Type) string {
	if message.Kind() == reflect.Pointer {
		return strcase.ToSnake(message.Elem().Name())
	}
	return strcase.ToSnake(message.Name())
}

func GetQueueName(message interface{}) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}
	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

func GetQueueNameFromType(message reflect.Type) string {
	if message.Kind() == reflect.Pointer {
		return strcase.ToSnake(message.Elem().Name())
	}
	return strcase.ToSnake(message.Name())
}

func GetRoutingKey(message interface{}) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}
	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

func GetRoutingKeyFromType(message reflect.Type) string {
	if message.Kind() == reflect.Pointer {
		return strcase.ToSnake(message.Elem().Name())
	}
	return strcase.ToSnake(message.Name())
}

func RegisterCustomMessageTypesToRegistrty(typesMap map[string]types.IMessage) {
	if typesMap == nil || len(typesMap) == 0 {
		return
	}

	for k, v := range typesMap {
		typeMapper.RegisterTypeWithKey(k, typeMapper.GetReflectType(v))
	}
}

func GetAllMessageTypes() []reflect.Type {
	var squares []reflect.Type
	d := linq.From(typeMapper.GetAllRegisteredTypes()).SelectManyT(func(i linq.KeyValue) linq.Query {
		return linq.From(i.Value)
	})
	d.ToSlice(&squares)
	res := typeMapper.TypesImplementedInterfaceWithFilterTypes[types.IMessage](squares)
	linq.From(res).Distinct().ToSlice(&squares)

	return squares
}
