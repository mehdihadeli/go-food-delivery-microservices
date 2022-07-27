package repositories

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/aggregate"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/read_models"
)

type ElasticOrderRepository interface {
	IndexOrder(ctx context.Context, order *read_models.OrderRead) error
	GetByID(ctx context.Context, orderId string) (*read_models.OrderRead, error)
	UpdateOrder(ctx context.Context, order *read_models.OrderRead) error
	Search(ctx context.Context, text string, pq *utils.ListQuery) (*utils.ListResult[*aggregate.Order], error)
}
