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
	dtosv1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/events/integrationevents"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/models"

	"github.com/mehdihadeli/go-mediatr"
)

type createProductHandler struct {
	fxparams.ProductHandlerParams
}

func NewCreateProductHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*CreateProduct, *dtos.CreateProductResponseDto] {
	return &createProductHandler{
		ProductHandlerParams: params,
	}
}

func (c *createProductHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*CreateProduct, *dtos.CreateProductResponseDto](
		c,
	)
}

func (c *createProductHandler) Handle(
	ctx context.Context,
	command *CreateProduct,
) (*dtos.CreateProductResponseDto, error) {
	product := &models.Product{
		ProductId:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	var createProductResult *dtos.CreateProductResponseDto

	err := c.Uow.Do(
		ctx,
		func(catalogContext contracts.CatalogContext) error {
			createdProduct, err := catalogContext.Products().
				CreateProduct(ctx, product)
			if err != nil {
				return customErrors.NewApplicationErrorWrapWithCode(
					err,
					http.StatusConflict,
					"product already exists",
				)
			}
			productDto, err := mapper.Map[*dtosv1.ProductDto](
				createdProduct,
			)
			if err != nil {
				return customErrors.NewApplicationErrorWrap(
					err,
					"error in the mapping ProductDto",
				)
			}

			productCreated := integrationevents.NewProductCreatedV1(
				productDto,
			)

			err = c.RabbitmqProducer.PublishMessage(ctx, productCreated, nil)
			if err != nil {
				return customErrors.NewApplicationErrorWrap(
					err,
					"error in publishing ProductCreated integration_events event",
				)
			}

			c.Log.Infow(
				fmt.Sprintf(
					"ProductCreated message with messageId `%s` published to the rabbitmq broker",
					productCreated.MessageId,
				),
				logger.Fields{"MessageId": productCreated.MessageId},
			)

			createProductResult = &dtos.CreateProductResponseDto{
				ProductID: product.ProductId,
			}

			c.Log.Infow(
				fmt.Sprintf(
					"product with id '%s' created",
					command.ProductID,
				),
				logger.Fields{
					"ProductId": command.ProductID,
					"MessageId": productCreated.MessageId,
				},
			)

			return nil
		},
	)

	return createProductResult, err
}
