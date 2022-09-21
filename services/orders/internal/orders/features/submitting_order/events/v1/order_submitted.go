package v1

import (
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	uuid "github.com/satori/go.uuid"
)

type OrderSubmittedV1 struct {
	OrderId uuid.UUID `json:"orderId" bson:"orderId,omitempty"`
}

func NewSubmitOrderV1(orderId uuid.UUID) (*OrderSubmittedV1, error) {
	if orderId == uuid.Nil {
		return nil, customErrors.NewDomainError(fmt.Sprintf("orderId {%s} is invalid", orderId))
	}

	event := OrderSubmittedV1{OrderId: orderId}

	return &event, nil
}
