package v1

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	httpErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/dto"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/shared/web"
	"net/http"
)

type CreateProductEndpoint struct {
	web.CatalogsEndpointBase
}

func NewCreteProductEndpoint(
	echo *echo.Echo,
	log logger.Logger,
	cfg *config.Config,
	mediator *mediatr.Mediator,
	validator *validator.Validate,
	metrics *shared.CatalogsServiceMetrics,
) *CreateProductEndpoint {
	return &CreateProductEndpoint{web.CatalogsEndpointBase{echo, log, cfg, mediator, validator, metrics}}
}

func (ep *CreateProductEndpoint) MapCreateProductEndpoint() {
	v1 := ep.Echo.Group("/api/v1")
	products := v1.Group("/" + ep.Cfg.Http.ProductsPath)
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
func (ep *CreateProductEndpoint) createProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		ep.Metrics.CreateProductHttpRequests.Inc()

		ctx, span := tracing.StartHttpServerTracerSpan(c, "productsHandlers.CreateProduct")
		defer span.Finish()

		request := &dto.CreateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Log.WarnMsg("Bind", err)
			ep.TraceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, ep.Cfg.Http.DebugErrorsResponse)
		}

		if err := ep.Validator.StructCtx(ctx, request); err != nil {
			ep.Log.Errorf("(validate) err: {%v}", err)
			tracing.TraceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, ep.Cfg.Http.DebugErrorsResponse)
		}

		command := creating_product.NewCreateProduct(request.Name, request.Description, request.Price)
		_, err := ep.Mediator.Send(ctx, *command)

		if err != nil {
			ep.Log.Errorf("(CreateOrder.Handle) id: {%s}, err: {%v}", command.ProductID, err)
			tracing.TraceErr(span, err)
			return httpErrors.ErrorCtxResponse(c, err, ep.Cfg.Http.DebugErrorsResponse)
		}

		ep.Log.Infof("(order created) id: {%s}", command.ProductID)
		ep.Metrics.SuccessHttpRequests.Inc()
		return c.JSON(http.StatusCreated, dto.CreateProductResponseDto{ProductID: command.ProductID})
	}
}
