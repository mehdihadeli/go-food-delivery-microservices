package commands

import (
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils/validator"
)

type CreateProduct struct {
	// we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
	Id          string    `validate:"required"`
	ProductId   string    `validate:"required"`
	Name        string    `validate:"required,min=3,max=250"`
	Description string    `validate:"required,min=3,max=500"`
	Price       float64   `validate:"required"`
	CreatedAt   time.Time `validate:"required"`
}

func NewCreateProduct(productId string, name string, description string, price float64, createdAt time.Time) (*CreateProduct, error) {
	command := &CreateProduct{Id: uuid.NewV4().String(), ProductId: productId, Name: name, Description: description, Price: price, CreatedAt: createdAt}
	err := validator.Validate(command)
	if err != nil {
		return nil, err
	}

	return command, nil
}
