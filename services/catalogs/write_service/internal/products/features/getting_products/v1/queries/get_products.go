package queries

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils/validator"
)

// Ref: https://golangbot.com/inheritance/

type GetProducts struct {
	*utils.ListQuery
}

func NewGetProducts(query *utils.ListQuery) (*GetProducts, error) {
	q := &GetProducts{ListQuery: query}

	err := validator.Validate(q)
	if err != nil {
		return nil, err
	}

	return q, nil
}
