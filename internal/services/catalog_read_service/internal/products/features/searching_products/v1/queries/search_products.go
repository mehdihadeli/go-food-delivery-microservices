package queries

import (
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
)

type SearchProducts struct {
	SearchText string
	*utils.ListQuery
}

func (s *SearchProducts) Validate() error {
	return validation.ValidateStruct(s, validation.Field(&s.SearchText, validation.Required))
}
