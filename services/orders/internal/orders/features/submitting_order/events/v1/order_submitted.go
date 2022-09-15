package v1

import (
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	uuid "github.com/satori/go.uuid"
)

type OrderSubmittedV1 struct {
	OrderID uuid.UUID `json:"orderID" bson:"orderID,omitempty"`
}

func NewSubmitOrderV1(orderID uuid.UUID) (*OrderSubmittedV1, error) {
	if orderID == uuid.Nil {
		return nil, customErrors.NewDomainError(fmt.Sprintf("orderId {%s} is invalid", orderID))
	}

	event := OrderSubmittedV1{OrderID: orderID}

	return &event, nil
}
