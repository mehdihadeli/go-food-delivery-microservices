package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/cqrs"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/mapper"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/gormdbcontext"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/data/datamodels"
	dto "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/updatingproduct/v1/events/integrationevents"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/models"

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

// IsTxRequest for enabling transactions on the mediatr pipeline
func (c *updateProductHandler) isTxRequest() {
}

func (c *updateProductHandler) Handle(
	ctx context.Context,
	command *UpdateProduct,
) (*mediatr.Unit, error) {
	product, err := gormdbcontext.FindModelByID[*datamodels.ProductDataModel, *models.Product](
		ctx,
		c.CatalogsDBContext,
		command.ProductID,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrapWithCode(
			err,
			http.StatusNotFound,
			fmt.Sprintf(
				"product with id `%s` not found",
				command.ProductID,
			),
		)
	}

	product.Name = command.Name
	product.Price = command.Price
	product.Description = command.Description
	product.UpdatedAt = command.UpdatedAt

	updatedProduct, err := gormdbcontext.UpdateModel[*datamodels.ProductDataModel, *models.Product](
		ctx,
		c.CatalogsDBContext,
		product,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in updating product in the repository",
		)
	}

	productDto, err := mapper.Map[*dto.ProductDto](updatedProduct)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping ProductDto",
		)
	}

	productUpdated := integrationevents.NewProductUpdatedV1(productDto)

	err = c.RabbitmqProducer.PublishMessage(ctx, productUpdated, nil)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in publishing 'ProductUpdated' message",
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"product with id '%s' updated",
			command.ProductID,
		),
		logger.Fields{"Id": command.ProductID},
	)

	c.Log.Infow(
		fmt.Sprintf(
			"ProductUpdated message with messageId `%s` published to the rabbitmq broker",
			productUpdated.MessageId,
		),
		logger.Fields{"MessageId": productUpdated.MessageId},
	)

	return &mediatr.Unit{}, err
}
