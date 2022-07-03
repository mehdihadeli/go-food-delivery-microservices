package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products"
	shared_configurations "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations"
	"net/http"
)

type getProductsEndpoint struct {
	mediator          *mediatr.Mediator
	productRepository repositories.ProductRepository
	productsGroup     *echo.Group
	infrastructure    *shared_configurations.Infrastructure
}

func NewGetProductsEndpoint(infra *shared_configurations.Infrastructure, mediator *mediatr.Mediator, productsGroup *echo.Group, productRepository repositories.ProductRepository) *getProductsEndpoint {
	return &getProductsEndpoint{mediator: mediator, productRepository: productRepository, productsGroup: productsGroup, infrastructure: infra}
}

func (ep *getProductsEndpoint) MapRoute() {
	ep.productsGroup.GET("", ep.getAllProducts())
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
		ep.infrastructure.Metrics.GetProductByIdHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "getProductsEndpoint.getAllProducts")
		defer span.Finish()

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, ep.infrastructure.Log, err)
			return err
		}

		request := &getting_products.GetProductsRequestDto{ListQuery: listQuery}
		if err := c.Bind(request); err != nil {
			ep.infrastructure.Log.WarnMsg("Bind", err)
			ep.infrastructure.TraceErr(span, err)
			return err
		}

		query := getting_products.GetProducts{request.ListQuery}

		queryResult, err := ep.mediator.Send(ctx, query)

		if err != nil {
			ep.infrastructure.Log.WarnMsg("GetProducts", err)
			ep.infrastructure.Metrics.ErrorHttpRequests.Inc()
			return err
		}

		response, ok := queryResult.(*getting_products.GetProductsResponseDto)
		err = utils.CheckType(ok)
		if err != nil {
			return err
		}

		ep.infrastructure.Metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusOK, response)
	}
}
