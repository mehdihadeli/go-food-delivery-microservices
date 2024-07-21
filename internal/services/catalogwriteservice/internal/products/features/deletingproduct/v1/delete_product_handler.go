package v1

import (
	"context"
	"fmt"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/cqrs"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/gormdbcontext"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/data/datamodels"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	integrationEvents "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/deletingproduct/v1/events/integrationevents"

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

// IsTxRequest for enabling transactions on the mediatr pipeline
func (c *deleteProductHandler) isTxRequest() {
}

func (c *deleteProductHandler) Handle(
	ctx context.Context,
	command *DeleteProduct,
) (*mediatr.Unit, error) {
	err := gormdbcontext.DeleteDataModelByID[*datamodels.ProductDataModel](ctx, c.CatalogsDBContext, command.ProductID)
	if err != nil {
		return nil, err
	}

	productDeleted := integrationEvents.NewProductDeletedV1(
		command.ProductID.String(),
	)

	if err = c.RabbitmqProducer.PublishMessage(ctx, productDeleted, nil); err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
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
		logger.Fields{"Id": command.ProductID},
	)

	return &mediatr.Unit{}, err
}
