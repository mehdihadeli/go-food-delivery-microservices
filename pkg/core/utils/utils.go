package utils

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"reflect"
)

func GetAllDomainEventTypes() []reflect.Type {
	var squares []reflect.Type
	d := linq.From(typeMapper.GetAllRegisteredTypes()).SelectManyT(func(i linq.KeyValue) linq.Query {
		return linq.From(i.Value)
	})
	d.ToSlice(&squares)
	res := typeMapper.TypesImplementedInterface[domain.IDomainEvent](squares)
	linq.From(res).Distinct().ToSlice(&squares)

	return squares
}

func GetAllEventTypes() []reflect.Type {
	var squares []reflect.Type
	d := linq.From(typeMapper.GetAllRegisteredTypes()).SelectManyT(func(i linq.KeyValue) linq.Query {
		return linq.From(i.Value)
	})
	d.ToSlice(&squares)
	res := typeMapper.TypesImplementedInterface[core.IEvent](squares)
	linq.From(res).Distinct().ToSlice(&squares)

	return squares
}
