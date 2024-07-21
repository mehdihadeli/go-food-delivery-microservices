package cqrs

import (
	"testing"
	"time"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/reflection/typemapper"

	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Command(t *testing.T) {
	command := &CreateProductTest{
		Command:     NewCommandByT[*CreateProductTest](),
		ProductID:   uuid.NewV4(),
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.AdjectiveDescriptive(),
		Price:       gofakeit.Price(100, 1000),
	}

	isImplementedCommand := typemapper.ImplementedInterfaceT[Command](command)
	assert.True(t, isImplementedCommand)

	var i interface{} = command
	_, isCommand := i.(Command)
	_, isTypeInfo := i.(TypeInfo)
	_, isQuery := i.(Query)
	_, isRequest := i.(Request)

	assert.True(t, isCommand)
	assert.True(t, isTypeInfo)
	assert.True(t, isRequest)
	assert.False(t, isQuery)

	assert.True(t, IsCommand(command))
	assert.True(t, IsRequest(command))
	assert.False(t, IsQuery(command))

	assert.Equal(t, command.ShortTypeName(), "*CreateProductTest")
	assert.Equal(t, command.FullTypeName(), "*cqrs.CreateProductTest")
}

type CreateProductTest struct {
	Command

	Name        string
	ProductID   uuid.UUID
	Description string
	Price       float64
	CreatedAt   time.Time
}
