package endpoints

import (
	"fmt"
	"net/http"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/web/route"
	customErrors "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/http/httperrors/customerrors"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/orders/contracts/params"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/orders/features/getting_order_by_id/v1/dtos"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/orders/features/getting_order_by_id/v1/queries"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type getOrderByIdEndpoint struct {
	params.OrderRouteParams
}

func NewGetOrderByIdEndpoint(params params.OrderRouteParams) route.Endpoint {
	return &getOrderByIdEndpoint{OrderRouteParams: params}
}

func (ep *getOrderByIdEndpoint) MapEndpoint() {
	ep.OrdersGroup.GET("/:id", ep.handler())
}

// Get Order By ID
// @Tags Orders
// @Summary Get order by id
// @Description Get order by id
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} dtos.GetOrderByIdResponseDto
// @Router /api/v1/orders/{id} [get]
func (ep *getOrderByIdEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		ep.OrdersMetrics.GetOrderByIdHttpRequests.Add(ctx, 1)

		request := &dtos.GetOrderByIdRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[getProductByIdEndpoint_handler.Bind] error in the binding request",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[getProductByIdEndpoint_handler.Bind] err: %v", badRequestErr),
			)
			return badRequestErr
		}

		query, err := queries.NewGetOrderById(request.Id)
		if err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"[getProductByIdEndpoint_handler.StructCtx]  query validation failed",
			)
			ep.Logger.Errorf("[getProductByIdEndpoint_handler.StructCtx] err: %v", validationErr)
			return validationErr
		}

		queryResult, err := mediatr.Send[*queries.GetOrderById, *dtos.GetOrderByIdResponseDto](
			ctx,
			query,
		)
		if err != nil {
			err = errors.WithMessage(
				err,
				"[getProductByIdEndpoint_handler.Send] error in sending GetOrderById",
			)
			ep.Logger.Errorw(
				fmt.Sprintf(
					"[getProductByIdEndpoint_handler.Send] id: {%s}, err: %v",
					query.Id,
					err,
				),
				logger.Fields{"Id": query.Id},
			)
			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
