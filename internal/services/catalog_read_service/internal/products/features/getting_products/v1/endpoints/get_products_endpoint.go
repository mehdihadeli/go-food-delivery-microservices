package endpoints

import (
	"fmt"
	"net/http"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/contracts/params"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web/route"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/getting_products/v1/dtos"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogreadservice/internal/products/features/getting_products/v1/queries"
)

type getProductsEndpoint struct {
	params.ProductRouteParams
}

func NewGetProductsEndpoint(
	params params.ProductRouteParams,
) route.Endpoint {
	return &getProductsEndpoint{
		ProductRouteParams: params,
	}
}

func (ep *getProductsEndpoint) MapEndpoint() {
	ep.ProductsGroup.GET("", ep.handler())
}

// GetAllProducts
// @Tags Products
// @Summary Get all product
// @Description Get all products
// @Accept json
// @Produce json
// @Param getProductsRequestDto query dtos.GetProductsRequestDto false "GetProductsRequestDto"
// @Success 200 {object} dtos.GetProductsResponseDto
// @Router /api/v1/products [get]
func (ep *getProductsEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		ep.CatalogsMetrics.GetProductsHttpRequests.Add(ctx, 1)

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[getProductsEndpoint_handler.GetListQueryFromCtx] error in getting data from query string",
			)
			ep.Logger.Errorf(
				fmt.Sprintf(
					"[getProductsEndpoint_handler.GetListQueryFromCtx] err: %v",
					badRequestErr,
				),
			)
			return err
		}

		request := queries.NewGetProducts(listQuery)
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[getProductsEndpoint_handler.Bind] error in the binding request",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[getProductsEndpoint_handler.Bind] err: %v", badRequestErr),
			)
			return badRequestErr
		}
		query := &queries.GetProducts{ListQuery: request.ListQuery}

		queryResult, err := mediatr.Send[*queries.GetProducts, *dtos.GetProductsResponseDto](
			ctx,
			query,
		)
		if err != nil {
			err = errors.WithMessage(
				err,
				"[getProductsEndpoint_handler.Send] error in sending GetProducts",
			)
			ep.Logger.Error(fmt.Sprintf("[getProductsEndpoint_handler.Send] err: {%v}", err))
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
