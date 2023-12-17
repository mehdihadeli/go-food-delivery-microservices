package v1

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
)

type CreateProduct struct {
	// we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
	Id          string
	ProductId   string
	Name        string
	Description string
	Price       float64
	CreatedAt   time.Time
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
	if err := command.Validate(); err != nil {
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
