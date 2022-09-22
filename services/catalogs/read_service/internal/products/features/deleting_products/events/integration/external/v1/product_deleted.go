package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
)

type ProductDeletedV1 struct {
	*types.Message
	ProductId string `json:"productId,omitempty"`
}
