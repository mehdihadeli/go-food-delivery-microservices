package v1

import (
	"github.com/labstack/echo/v4"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/dtos"
	shared_configurations "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/configurations"
	"net/http"
)

type createProductEndpoint struct {
	mediator          *mediatr.Mediator
	productRepository repositories.ProductRepository
	productsGroup     *echo.Group
	infrastructure    *shared_configurations.Infrastructure
}

func NewCreteProductEndpoint(infra *shared_configurations.Infrastructure, mediator *mediatr.Mediator, productsGroup *echo.Group, productRepository repositories.ProductRepository) *createProductEndpoint {
	return &createProductEndpoint{mediator: mediator, productRepository: productRepository, productsGroup: productsGroup, infrastructure: infra}
}

func (ep *createProductEndpoint) MapRoute() {
	ep.productsGroup.POST("", ep.createProduct())
}

// CreateProduct
// @Tags Products
// @Summary Create product
// @Description Create new product item
// @Accept json
// @Produce json
// @Param CreateProductRequestDto body dtos.CreateProductRequestDto true "Product data"
// @Success 201 {object} dtos.CreateProductResponseDto
// @Router /products [post]
func (ep *createProductEndpoint) createProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		ep.infrastructure.Metrics.CreateProductHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "createProductEndpoint.createProduct")
		defer span.Finish()

		request := &dtos.CreateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.infrastructure.Log.WarnMsg("Bind", err)
			ep.infrastructure.TraceErr(span, err)
			return httpErrors.ErrorResponse(err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		if err := ep.infrastructure.Validator.StructCtx(ctx, request); err != nil {
			ep.infrastructure.Log.Errorf("(validate) err: {%v}", err)
			tracing.TraceErr(span, err)
			return httpErrors.ErrorResponse(err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		command := creating_product.NewCreateProduct(request.Name, request.Description, request.Price)
		result, err := ep.mediator.Send(ctx, command)

		if err != nil {
			ep.infrastructure.Log.Errorf("(CreateOrder.Handle) id: {%s}, err: {%v}", command.ProductID, err)
			tracing.TraceErr(span, err)
			return httpErrors.ErrorResponse(err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		response, ok := result.(*dtos.CreateProductResponseDto)
		err = utils.CheckType(ok)
		if err != nil {
			return httpErrors.ErrorResponse(err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		ep.infrastructure.Log.Infof("(product created) id: {%s}", command.ProductID)
		ep.infrastructure.Metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusCreated, response)
	}
}
