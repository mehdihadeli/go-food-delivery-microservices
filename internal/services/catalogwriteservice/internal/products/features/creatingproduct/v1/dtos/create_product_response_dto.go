package dtos

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/serializer/json"

	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/response/
type CreateProductResponseDto struct {
	ProductID uuid.UUID `json:"productId"`
}

func (c *CreateProductResponseDto) String() string {
	return json.PrettyPrint(c)
}
