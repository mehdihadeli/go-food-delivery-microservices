package utils

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/iancoleman/strcase"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"reflect"
)

func GetMessageName(message interface{}) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}
	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

func GetTopicOrExchangeName(message interface{}) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}
	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

func GetQueueName(message interface{}) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}
	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

func GetRoutingKey(message interface{}) string {
	if reflect.TypeOf(message).Kind() == reflect.Pointer {
		return strcase.ToSnake(reflect.TypeOf(message).Elem().Name())
	}
	return strcase.ToSnake(reflect.TypeOf(message).Name())
}

func RegisterCustomMessageTypesToRegistrty(typesMap map[string]types.IMessage) {
	if typesMap == nil || len(typesMap) == 0 {
		return
	}

	for k, v := range typesMap {
		typeMapper.RegisterTypeWithKey(k, typeMapper.GetType(v))
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
