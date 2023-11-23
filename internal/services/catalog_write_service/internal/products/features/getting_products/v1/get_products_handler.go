package v1

import (
	"context"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/cqrs"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	dtosv1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_products/v1/dtos"

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
	products, err := c.ProductRepository.GetAllProducts(ctx, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtosv1.ProductDto](
		products,
	)
	if err != nil {
		return nil, err
	}

	c.Log.Info("products fetched")

	return &dtos.GetProductsResponseDto{Products: listResultDto}, nil
}
