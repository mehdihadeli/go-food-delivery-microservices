package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/cqrs"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	integrationEvents "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/deleting_product/v1/events/integrationevents"

	"github.com/mehdihadeli/go-mediatr"
)

type deleteProductHandler struct {
	fxparams.ProductHandlerParams
}

func NewDeleteProductHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*DeleteProduct, *mediatr.Unit] {
	return &deleteProductHandler{
		ProductHandlerParams: params,
	}
}

func (c *deleteProductHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*DeleteProduct, *mediatr.Unit](
		c,
	)
}

func (c *deleteProductHandler) Handle(
	ctx context.Context,
	command *DeleteProduct,
) (*mediatr.Unit, error) {
	err := c.Uow.Do(ctx, func(catalogContext contracts.CatalogContext) error {
		if err := catalogContext.Products().DeleteProductByID(ctx, command.ProductID); err != nil {
			return customErrors.NewApplicationErrorWrapWithCode(
				err,
				http.StatusNotFound,
				"product not found",
			)
		}

		productDeleted := integrationEvents.NewProductDeletedV1(
			command.ProductID.String(),
		)
		err := c.RabbitmqProducer.PublishMessage(ctx, productDeleted, nil)
		if err != nil {
			return customErrors.NewApplicationErrorWrap(
				err,
				"error in publishing 'ProductDeleted' message",
			)
		}

		c.Log.Infow(
			fmt.Sprintf(
				"ProductDeleted message with messageId '%s' published to the rabbitmq broker",
				productDeleted.MessageId,
			),
			logger.Fields{"MessageId": productDeleted.MessageId},
		)

		c.Log.Infow(
			fmt.Sprintf(
				"product with id '%s' deleted",
				command.ProductID,
			),
			logger.Fields{"ProductId": command.ProductID},
		)

		return nil
	})

	return &mediatr.Unit{}, err
}
