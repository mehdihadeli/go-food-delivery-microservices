package utils

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/domain"
	typeMapper "github.com/mehdihadeli/store-golang-microservice-sample/pkg/reflection/type_mappper"
	"reflect"
)

func GetAllDomainEventTypes() []reflect.Type {
	var types []reflect.Type
	d := linq.From(typeMapper.GetAllRegisteredTypes()).SelectManyT(func(i linq.KeyValue) linq.Query {
		return linq.From(i.Value)
	})
	d.ToSlice(&types)
	res := typeMapper.TypesImplementedInterfaceWithFilterTypes[domain.IDomainEvent](types)
	linq.From(res).Distinct().ToSlice(&types)

	return types
}

func GetAllEventTypes() []reflect.Type {
	var types []reflect.Type
	d := linq.From(typeMapper.GetAllRegisteredTypes()).SelectManyT(func(i linq.KeyValue) linq.Query {
		return linq.From(i.Value)
	})
	d.ToSlice(&types)
	res := typeMapper.TypesImplementedInterfaceWithFilterTypes[core.IEvent](types)
	linq.From(res).Distinct().ToSlice(&types)

	return types
}
