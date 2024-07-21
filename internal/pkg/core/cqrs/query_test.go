package cqrs

import (
	"testing"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/reflection/typemapper"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Query(t *testing.T) {
	query := &GetProductById{
		Query:     NewQueryByT[*GetProductById](),
		ProductID: uuid.NewV4(),
	}

	isImplementedQuery := typemapper.ImplementedInterfaceT[Query](query)
	assert.True(t, isImplementedQuery)

	var i interface{} = query
	_, isQuery := i.(Query)
	_, isTypeInfo := i.(TypeInfo)
	_, isCommand := i.(Command)
	_, isRequest := i.(Request)

	assert.True(t, isQuery)
	assert.False(t, isCommand)
	assert.True(t, isTypeInfo)
	assert.True(t, isRequest)

	assert.True(t, IsQuery(query))
	assert.False(t, IsCommand(query))
	assert.True(t, IsRequest(query))

	assert.Equal(t, query.ShortTypeName(), "*GetProductById")
	assert.Equal(t, query.FullTypeName(), "*cqrs.GetProductById")
}

type GetProductById struct {
	Query

	ProductID uuid.UUID
}
