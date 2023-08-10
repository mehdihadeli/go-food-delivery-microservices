package queries

import (
	"context"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/contracts/repositories"
	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/dtos/v1"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/features/getting_orders/v1/dtos"
)

type GetOrdersHandler struct {
	log                      logger.Logger
	mongoOrderReadRepository repositories.OrderMongoRepository
	tracer                   tracing.AppTracer
}

func NewGetOrdersHandler(
	log logger.Logger,
	mongoOrderReadRepository repositories.OrderMongoRepository,
	tracer tracing.AppTracer,
) *GetOrdersHandler {
	return &GetOrdersHandler{
		log:                      log,
		mongoOrderReadRepository: mongoOrderReadRepository,
		tracer:                   tracer,
	}
}

func (c *GetOrdersHandler) Handle(
	ctx context.Context,
	query *GetOrders,
) (*dtos.GetOrdersResponseDto, error) {
	ctx, span := c.tracer.Start(ctx, "GetOrdersHandler.Handle")
	span.SetAttributes(attribute.Object("Query", query))
	defer span.End()

	products, err := c.mongoOrderReadRepository.GetAllOrders(ctx, query.ListQuery)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				"[GetOrdersHandler_Handle.GetAllOrders] error in getting orders in the repository",
			),
		)
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtosV1.OrderReadDto](products)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(
			span,
			customErrors.NewApplicationErrorWrap(
				err,
				"[GetOrdersHandler_Handle.ListResultToListResultDto] error in the mapping ListResultToListResultDto",
			),
		)
	}

	c.log.Info("[GetOrdersHandler.Handle] orders fetched")

	return &dtos.GetOrdersResponseDto{Orders: listResultDto}, nil
}
