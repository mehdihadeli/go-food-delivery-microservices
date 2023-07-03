package endpoints

import (
	"fmt"
	"net/http"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/contracts/params"
	createProductCommand "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/commands"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/web/route"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/features/creating_product/v1/dtos"
)

type createProductEndpoint struct {
	params.ProductRouteParams
}

func NewCreteProductEndpoint(params params.ProductRouteParams) route.Endpoint {
	return &createProductEndpoint{ProductRouteParams: params}
}

func (ep *createProductEndpoint) MapEndpoint() {
	ep.ProductsGroup.POST("", ep.handler())
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
func (ep *createProductEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		ep.CatalogsMetrics.CreateProductHttpRequests.Add(ctx, 1)

		request := &dtos.CreateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[createProductEndpoint_handler.Bind] error in the binding request",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[createProductEndpoint_handler.Bind] err: %v", badRequestErr),
			)
		}

		command, err := createProductCommand.NewCreateProduct(
			request.Name,
			request.Description,
			request.Price,
		)
		if err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"[createProductEndpoint_handler.StructCtx] command validation failed",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[createProductEndpoint_handler.StructCtx] err: {%v}", validationErr),
			)
			return validationErr
		}

		result, err := mediatr.Send[*createProductCommand.CreateProduct, *dtos.CreateProductResponseDto](
			ctx,
			command,
		)
		if err != nil {
			err = errors.WithMessage(
				err,
				"[createProductEndpoint_handler.Send] error in sending CreateProduct",
			)
			ep.Logger.Errorw(
				fmt.Sprintf(
					"[createProductEndpoint_handler.Send] id: {%s}, err: {%v}",
					command.ProductID,
					err,
				),
				logger.Fields{"ProductId": command.ProductID},
			)
			return err
		}

		return c.JSON(http.StatusCreated, result)
	}
}
