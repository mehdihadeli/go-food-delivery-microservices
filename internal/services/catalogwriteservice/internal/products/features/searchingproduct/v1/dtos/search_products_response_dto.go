package dtos

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"
	dtoV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1"
)

type SearchProductsResponseDto struct {
	Products *utils.ListResult[*dtoV1.ProductDto]
}
