package pipeline

import (
	"context"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/validation"

	"github.com/mehdihadeli/go-mediatr"
)

type mediatorValidationPipeline struct {
	logger logger.Logger
}

func NewMediatorValidationPipeline(l logger.Logger) mediatr.PipelineBehavior {
	return &mediatorValidationPipeline{logger: l}
}

func (m mediatorValidationPipeline) Handle(
	ctx context.Context,
	request interface{},
	next mediatr.RequestHandlerFunc,
) (interface{}, error) {
	v, ok := request.(validation.Validator)
	if ok {
		err := v.Validate()
		if err != nil {
			validationErr := customErrors.NewValidationErrorWrap(err, "validation failed")

			return nil, validationErr
		}
	}

	return next(ctx)
}
