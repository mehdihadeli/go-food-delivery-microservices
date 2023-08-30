package commands

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"

	uuid "github.com/satori/go.uuid"
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

func NewCreateProduct(
	productId string,
	name string,
	description string,
	price float64,
	createdAt time.Time,
) (*CreateProduct, error) {
	command := &CreateProduct{
		Id:          uuid.NewV4().String(),
		ProductId:   productId,
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   createdAt,
	}
	err := command.Validate()
	if err != nil {
		return nil, err
	}

	return command, nil
}

func (p *CreateProduct) Validate() error {
	return validation.ValidateStruct(p, validation.Field(&p.Id, validation.Required),
		validation.Field(&p.ProductId, validation.Required),
		validation.Field(&p.Name, validation.Required, validation.Length(3, 250)),
		validation.Field(&p.Description, validation.Required, validation.Length(3, 500)),
		validation.Field(&p.Price, validation.Required),
		validation.Field(&p.CreatedAt, validation.Required))
}
