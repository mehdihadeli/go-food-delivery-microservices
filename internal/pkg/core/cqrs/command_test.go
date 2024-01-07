package cqrs

import (
	"testing"
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/reflection/typemapper"

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
	_, ok := i.(Command)
	_, ok2 := i.(TypeInfo)
	assert.True(t, ok)
	assert.True(t, ok2)
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
