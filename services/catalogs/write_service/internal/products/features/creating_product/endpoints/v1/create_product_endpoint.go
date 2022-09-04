package v1

import (
	"fmt"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product/commands/v1"
	"github.com/pkg/errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
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
	ep.ProductsGroup.POST("", ep.handler())
}

// CreateProductCommand
// @Tags Products
// @Summary Create product
// @Description Create new product item
// @Accept json
// @Produce json
// @Param CreateProductRequestDto body dtos.CreateProductRequestDto true "Product data"
// @Success 201 {object} dtos.CreateProductResponseDto
// @Router /api/v1/products [post]
func (ep *createProductEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ep.Metrics.CreateProductHttpRequests.Inc()
		ctx, span := tracing.StartHttpServerTracerSpan(c, "createProductEndpoint.handler")
		defer span.Finish()

		request := &dtos.CreateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[createProductEndpoint_handler.Bind] error in the binding request")
			ep.Log.Errorf(fmt.Sprintf("[createProductEndpoint_handler.Bind] err: %v", tracing.TraceWithErr(span, badRequestErr)))
			return badRequestErr
		}

		command := v1.NewCreateProductCommand(request.Name, request.Description, request.Price)
		if err := ep.Validator.StructCtx(ctx, command); err != nil {
			validationErr := customErrors.NewValidationErrorWrap(err, "[createProductEndpoint_handler.StructCtx] command validation failed")
			ep.Log.Errorf(fmt.Sprintf("[createProductEndpoint_handler.StructCtx] err: {%v}", tracing.TraceWithErr(span, validationErr)))
			return validationErr
		}

		result, err := mediatr.Send[*v1.CreateProductCommand, *dtos.CreateProductResponseDto](ctx, command)

		if err != nil {
			err = errors.WithMessage(err, "[createProductEndpoint_handler.Send] error in sending CreateProductCommand")
			ep.Log.Errorw(fmt.Sprintf("[createProductEndpoint_handler.Send] id: {%s}, err: {%v}", command.ProductID, tracing.TraceWithErr(span, err)), logger.Fields{"ProductId": command.ProductID})
			return err
		}

		return c.JSON(http.StatusCreated, result)
	}
}
