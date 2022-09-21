package external

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
)

type ProductDeleted struct {
	*types.Message
	ProductId string `json:"productId,omitempty"`
}
