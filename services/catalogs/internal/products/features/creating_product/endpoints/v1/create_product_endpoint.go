package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/contracts/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/dto"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/creating_product"
	shared_configurations "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/configurations"
	"net/http"
)

type createProductEndpoint struct {
	mediator          *mediatr.Mediator
	productRepository repositories.ProductRepository
	infrastructure    *shared_configurations.Infrastructure
}

func NewCreteProductEndpoint(infra *shared_configurations.Infrastructure, mediator *mediatr.Mediator, productRepository repositories.ProductRepository) *createProductEndpoint {
	return &createProductEndpoint{mediator: mediator, productRepository: productRepository, infrastructure: infra}
}

func (ep *createProductEndpoint) MapRoute() {
	v1 := ep.infrastructure.Echo.Group("/api/v1")
	products := v1.Group("/" + ep.infrastructure.Cfg.Http.ProductsPath)
	products.POST("", ep.createProduct())
}

// CreateProduct
// @Tags Products
// @Summary Create product
// @Description Create new product item
// @Accept json
// @Produce json
// @Param CreateProductRequestDto body dto.CreateProductRequestDto true "Product data"
// @Success 201 {object} dto.CreateProductResponseDto
// @Router /products [post]
func (ep *createProductEndpoint) createProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		ep.infrastructure.Metrics.CreateProductHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.CreateProduct")
		defer span.Finish()

		request := &dto.CreateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.infrastructure.Log.WarnMsg("Bind", err)
			ep.infrastructure.TraceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		if err := ep.infrastructure.Validator.StructCtx(ctx, request); err != nil {
			ep.infrastructure.Log.Errorf("(validate) err: {%v}", err)
			tracing.TraceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		command := creating_product.NewCreateProduct(request.Name, request.Description, request.Price)
		_, err := ep.mediator.Send(ctx, *command)

		if err != nil {
			ep.infrastructure.Log.Errorf("(CreateOrder.Handle) id: {%s}, err: {%v}", command.ProductID, err)
			tracing.TraceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, ep.infrastructure.Cfg.Http.DebugErrorsResponse)
		}

		ep.infrastructure.Log.Infof("(order created) id: {%s}", command.ProductID)
		ep.infrastructure.Metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusCreated, dto.CreateProductResponseDto{ProductID: command.ProductID})
	}
}
