package queries

import (
	validation "github.com/go-ozzo/ozzo-validation"
	uuid "github.com/satori/go.uuid"
)

type GetOrderById struct {
	Id uuid.UUID
}

func NewGetOrderById(id uuid.UUID) (*GetOrderById, error) {
	query := &GetOrderById{Id: id}

	err := query.Validate()
	if err != nil {
		return nil, err
	}

	return query, nil
}

func (g GetOrderById) Validate() error {
	return validation.ValidateStruct(&g,
		validation.Field(&g.Id, validation.Required),
	)
}
