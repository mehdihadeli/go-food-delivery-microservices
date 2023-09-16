package utils

import (
	"reflect"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/domain"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/events"
	typeMapper "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/type_mappper"

	"github.com/ahmetb/go-linq/v3"
)

func GetAllDomainEventTypes() []reflect.Type {
	var types []reflect.Type
	d := linq.From(typeMapper.GetAllRegisteredTypes()).
		SelectManyT(func(i linq.KeyValue) linq.Query {
			return linq.From(i.Value)
		})
	d.ToSlice(&types)
	res := typeMapper.TypesImplementedInterfaceWithFilterTypes[domain.IDomainEvent](types)
	linq.From(res).Distinct().ToSlice(&types)

	return types
}

func GetAllEventTypes() []reflect.Type {
	var types []reflect.Type
	d := linq.From(typeMapper.GetAllRegisteredTypes()).
		SelectManyT(func(i linq.KeyValue) linq.Query {
			return linq.From(i.Value)
		})
	d.ToSlice(&types)
	res := typeMapper.TypesImplementedInterfaceWithFilterTypes[events.IEvent](types)
	linq.From(res).Distinct().ToSlice(&types)

	return types
}
