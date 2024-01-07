package cqrs

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Query(t *testing.T) {
	query := &GetProductById{
		Query:     NewQueryByT[GetProductById](),
		ProductID: uuid.NewV4(),
	}

	assert.True(t, IsQuery(query))
}

//
//func Test_Query_Is_Catstable_To_Command(t *testing.T) {
//	var q Query = NewQuery()
//	var c Command = commands.NewCommand()
//	query, qok := q.(Query)
//	command, cok := c.(Command)
//	assert.True(t, qok)
//	assert.True(t, cok)
//	assert.NotNil(t, query)
//	assert.NotNil(t, command)
//
//	query, qok = command.(Query)
//	assert.False(t, qok)
//	assert.Nil(t, query)
//}

type GetProductById struct {
	Query
	ProductID uuid.UUID
}
