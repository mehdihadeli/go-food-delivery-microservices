package cqrs

type query struct {
	TypeInfo
}

type Query interface {
	TypeInfo
	isQuery()
}

func NewQueryByT[T any]() Query {
	return &query{TypeInfo: NewTypeInfoT[T]()}
}

func (q *query) isQuery() {
}

func IsQuery(obj interface{}) bool {
	if _, ok := obj.(Query); ok {
		return true
	}

	return false
}
