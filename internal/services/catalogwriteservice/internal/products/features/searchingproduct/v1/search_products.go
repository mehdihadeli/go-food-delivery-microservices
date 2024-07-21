package v1

import (
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"

	validation "github.com/go-ozzo/ozzo-validation"
)

type SearchProducts struct {
	SearchText string
	*utils.ListQuery
}

func NewSearchProducts(searchText string, query *utils.ListQuery) *SearchProducts {
	searchProductQuery := &SearchProducts{
		SearchText: searchText,
		ListQuery:  query,
	}

	return searchProductQuery
}

func NewSearchProductsWithValidation(searchText string, query *utils.ListQuery) (*SearchProducts, error) {
	searchProductQuery := NewSearchProducts(searchText, query)

	err := searchProductQuery.Validate()

	return searchProductQuery, err
}

func (p *SearchProducts) Validate() error {
	err := validation.ValidateStruct(p, validation.Field(&p.SearchText, validation.Required))
	if err != nil {
		return customErrors.NewValidationErrorWrap(err, "validation error")
	}

	return nil
}
