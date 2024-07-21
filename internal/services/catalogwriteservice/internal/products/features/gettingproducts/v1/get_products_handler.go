package v1

import (
	"context"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/cqrs"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/postgresgorm/helpers/gormextensions"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/utils"
	datamodel "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/data/datamodels"
	dtosv1 "github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/features/gettingproducts/v1/dtos"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/catalogwriteservice/internal/products/models"

	"github.com/mehdihadeli/go-mediatr"
)

type getProductsHandler struct {
	fxparams.ProductHandlerParams
}

func NewGetProductsHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*GetProducts, *dtos.GetProductsResponseDto] {
	return &getProductsHandler{
		ProductHandlerParams: params,
	}
}

func (c *getProductsHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*GetProducts, *dtos.GetProductsResponseDto](
		c,
	)
}

func (c *getProductsHandler) Handle(
	ctx context.Context,
	query *GetProducts,
) (*dtos.GetProductsResponseDto, error) {
	products, err := gormextensions.Paginate[*datamodel.ProductDataModel, *models.Product](
		ctx,
		query.ListQuery,
		c.CatalogsDBContext.DB(),
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the fetching products",
		)
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtosv1.ProductDto](
		products,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping",
		)
	}

	c.Log.Info("products fetched")

	return &dtos.GetProductsResponseDto{Products: listResultDto}, nil
}
