package repositories

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/orderservice/internal/orders/models/orders/read_models"
)

type orderReadRepository interface {
	GetAllOrders(
		ctx context.Context,
		listQuery *utils.ListQuery,
	) (*utils.ListResult[*read_models.OrderReadModel], error)
	SearchOrders(
		ctx context.Context,
		searchText string,
		listQuery *utils.ListQuery,
	) (*utils.ListResult[*read_models.OrderReadModel], error)
	GetOrderById(ctx context.Context, uuid uuid.UUID) (*read_models.OrderReadModel, error)
	GetOrderByOrderId(ctx context.Context, orderId uuid.UUID) (*read_models.OrderReadModel, error)
	CreateOrder(
		ctx context.Context,
		order *read_models.OrderReadModel,
	) (*read_models.OrderReadModel, error)
	UpdateOrder(
		ctx context.Context,
		order *read_models.OrderReadModel,
	) (*read_models.OrderReadModel, error)
	DeleteOrderByID(ctx context.Context, uuid uuid.UUID) error
}

type OrderElasticRepository interface {
	orderReadRepository
}

type OrderMongoRepository interface {
	orderReadRepository
}
