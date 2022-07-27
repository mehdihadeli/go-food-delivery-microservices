package v1

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/commands/v1"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/delivery"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/dtos"
)

type createProductEndpoint struct {
	*delivery.ProductEndpointBase
}

func NewCreteProductEndpoint(endpointBase *delivery.ProductEndpointBase) *createProductEndpoint {
	return &createProductEndpoint{endpointBase}
}

func (ep *createProductEndpoint) MapRoute() {
	ep.ProductsGroup.POST("", ep.createProduct())
}

// CreateProduct
// @Tags Products
// @Summary Create product
// @Description Create new product item
// @Accept json
// @Produce json
// @Param CreateProductRequestDto body dtos.CreateProductRequestDto true "Product data"
// @Success 201 {object} dtos.CreateProductResponseDto
// @Router /api/v1/products [post]
func (ep *createProductEndpoint) createProduct() echo.HandlerFunc {
	return func(c echo.Context) error {

		ep.Metrics.CreateProductHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "createProductEndpoint.createProduct")
		defer span.Finish()

		request := &dtos.CreateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			ep.Log.WarnMsg("Bind", err)
			tracing.TraceErr(span, err)
			return err
		}

		if err := ep.Validator.StructCtx(ctx, request); err != nil {
			ep.Log.Errorf("(validate) err: {%v}", err)
			tracing.TraceErr(span, err)
			return err
		}

		command := v1.NewCreateProduct(request.Name, request.Description, request.Price)
		result, err := mediatr.Send[*dtos.CreateProductResponseDto](ctx, command)

		if err != nil {
			ep.Log.Errorf("(CreateProduct.Handle) id: {%s}, err: {%v}", command.ProductID, err)
			tracing.TraceErr(span, err)
			return err
		}

		ep.Log.Infof("(product created) id: {%s}", command.ProductID)
		return c.JSON(http.StatusCreated, result)
	}
}
