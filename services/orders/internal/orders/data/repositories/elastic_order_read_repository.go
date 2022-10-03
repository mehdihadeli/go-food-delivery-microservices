package repositories

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/utils"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/read_models"
	uuid "github.com/satori/go.uuid"
)

type elasticOrderReadRepository struct {
	log           logger.Logger
	cfg           *config.Config
	elasticClient *elasticsearch.Client
}

func NewElasticOrderReadRepository(log logger.Logger, cfg *config.Config, elasticClient *elasticsearch.Client) repositories.OrderReadRepository {
	return &elasticOrderReadRepository{log: log, cfg: cfg, elasticClient: elasticClient}
}

func (e elasticOrderReadRepository) GetAllOrders(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*read_models.OrderReadModel], error) {
	//TODO implement me
	panic("implement me")
}

func (e elasticOrderReadRepository) SearchOrders(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[*read_models.OrderReadModel], error) {
	//TODO implement me
	panic("implement me")
}

func (e elasticOrderReadRepository) GetOrderById(ctx context.Context, uuid uuid.UUID) (*read_models.OrderReadModel, error) {
	//TODO implement me
	panic("implement me")
}

func (e elasticOrderReadRepository) GetOrderByOrderId(ctx context.Context, uuid uuid.UUID) (*read_models.OrderReadModel, error) {
	//TODO implement me
	panic("implement me")
}

func (e elasticOrderReadRepository) CreateOrder(ctx context.Context, order *read_models.OrderReadModel) (*read_models.OrderReadModel, error) {
	//TODO implement me
	panic("implement me")
}

func (e elasticOrderReadRepository) UpdateOrder(ctx context.Context, order *read_models.OrderReadModel) (*read_models.OrderReadModel, error) {
	//TODO implement me
	panic("implement me")
}

func (e elasticOrderReadRepository) DeleteOrderByID(ctx context.Context, uuid uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
