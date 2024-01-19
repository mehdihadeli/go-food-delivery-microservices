package cqrs

type query struct {
	TypeInfo
	Request
}

type Query interface {
	isQuery()

	Request
	TypeInfo
}

func NewQueryByT[T any]() Query {
	return &query{
		TypeInfo: NewTypeInfoT[T](),
		Request:  NewRequest(),
	}
}

func (q *query) isQuery() {
}

func IsQuery(obj interface{}) bool {
	if _, ok := obj.(Query); ok {
		return true
	}

	return false
}
