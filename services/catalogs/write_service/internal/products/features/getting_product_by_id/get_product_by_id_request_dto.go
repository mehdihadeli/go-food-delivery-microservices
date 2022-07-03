package getting_product_by_id

import uuid "github.com/satori/go.uuid"

type GetProductByIdRequestDto struct {
	ProductId uuid.UUID `param:"id" json:"-" validate:"required"`
}
