package queries

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"

	validation "github.com/go-ozzo/ozzo-validation"
)

type SearchProducts struct {
	SearchText string
	*utils.ListQuery
}

func (s *SearchProducts) Validate() error {
	return validation.ValidateStruct(s, validation.Field(&s.SearchText, validation.Required))
}
