package cqrs

import (
	"reflect"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/reflection/typemapper"
)

type TypeInfo interface {
	ShortTypeName() string
	FullTypeName() string
	Type() reflect.Type
}

type typeInfo struct {
	shortTypeName string
	fullTypeName  string
	typ           reflect.Type
}

func NewTypeInfoT[T any]() TypeInfo {
	name := typemapper.GetGenericTypeNameByT[T]()
	fullName := typemapper.GetGenericFullTypeNameByT[T]()
	typ := typemapper.GetGenericTypeByT[T]()

	return &typeInfo{fullTypeName: fullName, typ: typ, shortTypeName: name}
}

func (t *typeInfo) ShortTypeName() string {
	return t.shortTypeName
}

func (t *typeInfo) FullTypeName() string {
	return t.fullTypeName
}

func (t *typeInfo) Type() reflect.Type {
	return t.typ
}
