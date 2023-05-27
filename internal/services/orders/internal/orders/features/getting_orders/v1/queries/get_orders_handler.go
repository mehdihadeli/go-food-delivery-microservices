package queries

import (
	"context"

	customErrors "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing/attribute"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"
	dtosV1 "github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/dtos/v1"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/contracts/repositories"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orders/internal/orders/features/getting_orders/v1/dtos"
)

type GetOrdersHandler struct {
	log                      logger.Logger
	cfg                      *config.Config
	mongoOrderReadRepository repositories.OrderReadRepository
}

func NewGetOrdersHandler(log logger.Logger, cfg *config.Config, mongoOrderReadRepository repositories.OrderReadRepository) *GetOrdersHandler {
	return &GetOrdersHandler{log: log, cfg: cfg, mongoOrderReadRepository: mongoOrderReadRepository}
}

func (c *GetOrdersHandler) Handle(ctx context.Context, query *GetOrders) (*dtos.GetOrdersResponseDto, error) {
	ctx, span := tracing.Tracer.Start(ctx, "GetOrdersHandler.Handle")
	span.SetAttributes(attribute.Object("Query", query))
	defer span.End()

	products, err := c.mongoOrderReadRepository.GetAllOrders(ctx, query.ListQuery)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, "[GetOrdersHandler_Handle.GetAllOrders] error in getting orders in the repository"))
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtosV1.OrderReadDto](products)
	if err != nil {
		return nil, tracing.TraceErrFromSpan(span, customErrors.NewApplicationErrorWrap(err, "[GetOrdersHandler_Handle.ListResultToListResultDto] error in the mapping ListResultToListResultDto"))
	}

	c.log.Info("[GetOrdersHandler.Handle] orders fetched")

	return &dtos.GetOrdersResponseDto{Orders: listResultDto}, nil
}
