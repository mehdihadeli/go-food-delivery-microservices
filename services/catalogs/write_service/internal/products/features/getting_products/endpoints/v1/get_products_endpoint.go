package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/queries/v1"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products/dtos"
)

type getProductsEndpoint struct {
	*delivery.ProductEndpointBase
}

func NewGetProductsEndpoint(productEndpointBase *delivery.ProductEndpointBase) *getProductsEndpoint {
	return &getProductsEndpoint{productEndpointBase}
}

func (ep *getProductsEndpoint) MapRoute() {
	ep.ProductsGroup.GET("", ep.getAllProducts())
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
func (ep *getProductsEndpoint) getAllProducts() echo.HandlerFunc {
	return func(c echo.Context) error {

		ep.Metrics.GetProductsHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "getProductsEndpoint.getAllProducts")
		defer span.Finish()

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			return err
		}

		request := &dtos.GetProductsRequestDto{ListQuery: listQuery}
		if err := c.Bind(request); err != nil {
			ep.Log.WarnMsg("Bind", err)
			tracing.TraceErr(span, err)
			return err
		}

		query := &v1.GetProductsQuery{ListQuery: request.ListQuery}

		queryResult, err := mediatr.Send[*v1.GetProductsQuery, *dtos.GetProductsResponseDto](ctx, query)

		if err != nil {
			ep.Log.WarnMsg("GetProductsQuery", err)
			tracing.TraceErr(span, err)
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
