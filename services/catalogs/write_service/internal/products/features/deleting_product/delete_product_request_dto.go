package deleting_product

import uuid "github.com/satori/go.uuid"

type DeleteProductRequestDto struct {
	ProductID uuid.UUID `param:"id"  json:"-"`
}
