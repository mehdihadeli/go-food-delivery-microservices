package mediator

import "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/cqrs"

func RegisterMediatorHandlers(handlers []cqrs.HandlerRegisterer) error {
	for _, handler := range handlers {
		err := handler.RegisterHandler()
		if err != nil {
			return err
		}
	}

	return nil
}
