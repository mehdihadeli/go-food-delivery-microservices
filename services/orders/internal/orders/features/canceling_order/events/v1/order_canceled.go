package v1

type OrderCanceledEventV1 struct {
	CancelReason string `json:"cancelReason"`
}

func NewOrderCanceledEventV1(cancelReason string) (*OrderCanceledEventV1, error) {
	eventData := OrderCanceledEventV1{CancelReason: cancelReason}

	return &eventData, nil
}
