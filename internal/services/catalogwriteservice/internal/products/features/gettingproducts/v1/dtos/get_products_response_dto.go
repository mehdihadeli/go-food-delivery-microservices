package dtos

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"
	dtoV1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1"
)

// https://echo.labstack.com/guide/response/
type GetProductsResponseDto struct {
	Products *utils.ListResult[*dtoV1.ProductDto]
}
