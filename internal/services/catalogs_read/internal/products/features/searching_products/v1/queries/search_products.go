package queries

import "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

type SearchProducts struct {
	SearchText string `validate:"required"`
	*utils.ListQuery
}
