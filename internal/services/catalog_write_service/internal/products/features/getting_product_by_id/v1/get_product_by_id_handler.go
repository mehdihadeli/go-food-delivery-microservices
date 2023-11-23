package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/cqrs"
	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/mapper"
	dtoV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/getting_product_by_id/v1/dtos"

	"github.com/mehdihadeli/go-mediatr"
)

type getProductByIDHandler struct {
	fxparams.ProductHandlerParams
}

func NewGetProductByIDHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*GetProductById, *dtos.GetProductByIdResponseDto] {
	return &getProductByIDHandler{
		ProductHandlerParams: params,
	}
}

func (c *getProductByIDHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*GetProductById, *dtos.GetProductByIdResponseDto](
		c,
	)
}

func (c *getProductByIDHandler) Handle(
	ctx context.Context,
	query *GetProductById,
) (*dtos.GetProductByIdResponseDto, error) {
	product, err := c.ProductRepository.GetProductById(ctx, query.ProductID)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrapWithCode(
			err,
			http.StatusNotFound,
			fmt.Sprintf(
				"error in getting product with id %s in the repository",
				query.ProductID.String(),
			),
		)
	}

	productDto, err := mapper.Map[*dtoV1.ProductDto](product)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping product",
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"product with id: {%s} fetched",
			query.ProductID,
		),
		logger.Fields{"ProductId": query.ProductID.String()},
	)

	return &dtos.GetProductByIdResponseDto{Product: productDto}, nil
}
