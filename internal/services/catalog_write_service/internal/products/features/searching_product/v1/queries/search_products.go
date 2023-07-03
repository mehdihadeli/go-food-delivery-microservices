package queries

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/validator"
)

type SearchProducts struct {
	SearchText string `validate:"required"`
	*utils.ListQuery
}

func NewSearchProducts(searchText string, query *utils.ListQuery) (*SearchProducts, error) {
	command := &SearchProducts{
		SearchText: searchText,
		ListQuery:  query,
	}

	err := validator.Validate(command)
	if err != nil {
		return nil, err
	}

	return command, nil
}
