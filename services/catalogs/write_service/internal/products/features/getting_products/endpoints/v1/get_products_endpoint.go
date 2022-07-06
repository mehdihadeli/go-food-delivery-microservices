package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products"
	"net/http"
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
// @Param getProductsRequestDto query getting_products.GetProductsRequestDto false "GetProductsRequestDto"
// @Success 200 {object} getting_products.GetProductsResponseDto
// @Router /products [get]
func (ep *getProductsEndpoint) getAllProducts() echo.HandlerFunc {
	return func(c echo.Context) error {

		ep.Metrics.GetProductsHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "getProductsEndpoint.getAllProducts")
		defer span.Finish()

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, ep.Log, err)
			return err
		}

		request := &getting_products.GetProductsRequestDto{ListQuery: listQuery}
		if err := c.Bind(request); err != nil {
			ep.Log.WarnMsg("Bind", err)
			tracing.TraceErr(span, err)
			return err
		}

		query := getting_products.GetProducts{request.ListQuery}

		queryResult, err := ep.ProductMediator.Send(ctx, query)

		if err != nil {
			ep.Log.WarnMsg("GetProducts", err)
			tracing.TraceErr(span, err)
			return err
		}

		response, ok := queryResult.(*getting_products.GetProductsResponseDto)
		err = utils.CheckType(ok)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		return c.JSON(http.StatusOK, response)
	}
}
