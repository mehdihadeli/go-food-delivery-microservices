package dtos

import dtoV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dto/v1"

// https://echo.labstack.com/guide/response/
type GetProductByIdResponseDto struct {
	Product *dtoV1.ProductDto `json:"product"`
}
