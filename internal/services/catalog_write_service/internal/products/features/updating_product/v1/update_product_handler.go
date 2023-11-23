package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/cqrs"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts"
	dto "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/updating_product/v1/events/integrationevents"

	"github.com/mehdihadeli/go-mediatr"
)

type updateProductHandler struct {
	fxparams.ProductHandlerParams
	cqrs.HandlerRegisterer
}

func NewUpdateProductHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*UpdateProduct, *mediatr.Unit] {
	return &updateProductHandler{
		ProductHandlerParams: params,
	}
}

func (c *updateProductHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*UpdateProduct, *mediatr.Unit](
		c,
	)
}

func (c *updateProductHandler) Handle(
	ctx context.Context,
	command *UpdateProduct,
) (*mediatr.Unit, error) {
	err := c.Uow.Do(ctx, func(catalogContext contracts.CatalogContext) error {
		product, err := catalogContext.Products().
			GetProductById(ctx, command.ProductID)
		if err != nil {
			return customErrors.NewApplicationErrorWrapWithCode(
				err,
				http.StatusNotFound,
				fmt.Sprintf(
					"product with id %s not found",
					command.ProductID,
				),
			)
		}

		product.Name = command.Name
		product.Price = command.Price
		product.Description = command.Description
		product.UpdatedAt = command.UpdatedAt

		updatedProduct, err := catalogContext.Products().
			UpdateProduct(ctx, product)
		if err != nil {
			return customErrors.NewApplicationErrorWrap(
				err,
				"error in updating product in the repository",
			)
		}

		productDto, err := mapper.Map[*dto.ProductDto](updatedProduct)
		if err != nil {
			return customErrors.NewApplicationErrorWrap(
				err,
				"error in the mapping ProductDto",
			)
		}

		productUpdated := integrationevents.NewProductUpdatedV1(productDto)

		err = c.RabbitmqProducer.PublishMessage(ctx, productUpdated, nil)
		if err != nil {
			return customErrors.NewApplicationErrorWrap(
				err,
				"error in publishing 'ProductUpdated' message",
			)
		}

		c.Log.Infow(
			fmt.Sprintf(
				"product with id '%s' updated",
				command.ProductID,
			),
			logger.Fields{"ProductId": command.ProductID},
		)

		c.Log.Infow(
			fmt.Sprintf(
				"ProductUpdated message with messageId `%s` published to the rabbitmq broker",
				productUpdated.MessageId,
			),
			logger.Fields{"MessageId": productUpdated.MessageId},
		)

		return nil
	})

	return &mediatr.Unit{}, err
}
