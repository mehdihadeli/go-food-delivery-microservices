package queries

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/validator"
)

// Ref: https://golangbot.com/inheritance/

type GetProducts struct {
	*utils.ListQuery
}

func NewGetProducts(query *utils.ListQuery) (*GetProducts, error) {
	q := &GetProducts{ListQuery: query}

	// TODO
	// since there is no validate tag on ListQuery,
	// maybe we should just remove the next line ?
	err := validator.Validate(q)
	if err != nil {
		return nil, err
	}

	return q, nil
}
