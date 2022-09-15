package repositories

import (
	"context"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/read_models"
	uuid "github.com/satori/go.uuid"
)

type OrderReadRepository interface {
	GetAllOrders(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*read_models.OrderReadModel], error)
	SearchOrders(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[*read_models.OrderReadModel], error)
	GetOrderById(ctx context.Context, uuid uuid.UUID) (*read_models.OrderReadModel, error)
	CreateOrder(ctx context.Context, order *read_models.OrderReadModel) (*read_models.OrderReadModel, error)
	UpdateOrder(ctx context.Context, order *read_models.OrderReadModel) (*read_models.OrderReadModel, error)
	DeleteOrderByID(ctx context.Context, uuid uuid.UUID) error
}
