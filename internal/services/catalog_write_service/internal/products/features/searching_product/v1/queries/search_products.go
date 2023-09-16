package queries

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

	validation "github.com/go-ozzo/ozzo-validation"
)

type SearchProducts struct {
	SearchText string
	*utils.ListQuery
}

func NewSearchProducts(searchText string, query *utils.ListQuery) (*SearchProducts, error) {
	command := &SearchProducts{
		SearchText: searchText,
		ListQuery:  query,
	}

	err := command.Validate()
	if err != nil {
		return nil, err
	}

	return command, nil
}

func (p *SearchProducts) Validate() error {
	return validation.ValidateStruct(p, validation.Field(&p.SearchText, validation.Required))
}
