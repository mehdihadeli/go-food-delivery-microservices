package domain

import (
	"fmt"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
)

var (
	ErrCancelReasonRequired           = httpErrors.NewDomainError("Cancel reason must be provided")
	ErrOrderMustBePaidBeforeDelivered = httpErrors.NewDomainError("Order must be paid before been delivered")
	ErrOrderShopItemsIsRequired       = httpErrors.NewDomainError("order shop items is required")
	ErrInvalidDeliveryAddress         = httpErrors.NewDomainError("Invalid delivery address")
	ErrInvalidDeliveryTimeStamp       = httpErrors.NewDomainError("Invalid delivery timestamp")
	ErrInvalidAccountEmail            = httpErrors.NewDomainError("Invalid account email")
	ErrInvalidOrderID                 = httpErrors.NewDomainError("Invalid order id")
	ErrInvalidTime                    = httpErrors.NewDomainError("Invalid time")

	ErrOrderAlreadyCompleted = func(id int) httpErrors.ProblemDetailErr {
		return httpErrors.NewDomainError(fmt.Sprintf("Order with id '%v' already completed", id))
	}
	ErrOrderAlreadyCanceled = func(id int) httpErrors.ProblemDetailErr {
		return httpErrors.NewDomainError(fmt.Sprintf("Order with id '%v' already canceled", id))
	}
	ErrOrderAlreadyCancelled = func(id int) httpErrors.ProblemDetailErr {
		return httpErrors.NewDomainError(fmt.Sprintf("Order with id '%v' already cancelled", id))
	}
	ErrAlreadyPaid = func(id int) httpErrors.ProblemDetailErr {
		return httpErrors.NewDomainError(fmt.Sprintf("Order with id '%v' already paid", id))
	}
	ErrAlreadySubmitted = func(id int) httpErrors.ProblemDetailErr {
		return httpErrors.NewDomainError(fmt.Sprintf("Order with id '%v' already submitted", id))
	}
	ErrOrderNotPaid = func(id int) httpErrors.ProblemDetailErr {
		return httpErrors.NewDomainError(fmt.Sprintf("Order with id '%v' not paid", id))
	}
	ErrOrderNotFound = func(id int) httpErrors.ProblemDetailErr {
		return httpErrors.NewNotFoundError(fmt.Sprintf("Order with id '%v' not found", id))
	}
	ErrAlreadyCreated = func(id int) httpErrors.ProblemDetailErr {
		return httpErrors.NewDomainError(fmt.Sprintf("Order with id '%v' already created", id))
	}
)
