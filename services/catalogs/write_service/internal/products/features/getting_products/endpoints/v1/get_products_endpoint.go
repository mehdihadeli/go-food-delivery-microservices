package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/dtos"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/queries/v1"
	"github.com/pkg/errors"
	"net/http"
)

type getProductsEndpoint struct {
	*delivery.ProductEndpointBase
}

func NewGetProductsEndpoint(productEndpointBase *delivery.ProductEndpointBase) *getProductsEndpoint {
	return &getProductsEndpoint{productEndpointBase}
}

func (ep *getProductsEndpoint) MapRoute() {
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
		ep.Metrics.GetProductsHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "getProductsEndpoint.handler")
		defer span.Finish()

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[getProductsEndpoint_handler.GetListQueryFromCtx] error in getting data from query string")
			ep.Log.Errorf(fmt.Sprintf("[getProductsEndpoint_handler.GetListQueryFromCtx] err: %v", tracing.TraceWithErr(span, badRequestErr)))
			return err
		}

		request := &dtos.GetProductsRequestDto{ListQuery: listQuery}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[getProductsEndpoint_handler.Bind] error in the binding request")
			ep.Log.Errorf(fmt.Sprintf("[getProductsEndpoint_handler.Bind] err: %v", tracing.TraceWithErr(span, badRequestErr)))
			return badRequestErr
		}

		query := v1.NewGetProductsQuery(request.ListQuery)

		queryResult, err := mediatr.Send[*v1.GetProductsQuery, *dtos.GetProductsResponseDto](ctx, query)

		if err != nil {
			err = errors.WithMessage(err, "[getProductsEndpoint_handler.Send] error in sending GetProductsQuery")
			ep.Log.Error(fmt.Sprintf("[getProductsEndpoint_handler.Send] err: {%v}", tracing.TraceWithErr(span, err)))
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
