package creating_product

import (
	"context"
	"encoding/json"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/dtos"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CreateProductHandler struct {
	log           logger.Logger
	cfg           *config.Config
	repository    contracts.ProductRepository
	kafkaProducer kafkaClient.Producer
}

func NewCreateProductHandler(log logger.Logger, cfg *config.Config, repository contracts.ProductRepository, kafkaProducer kafkaClient.Producer) *CreateProductHandler {
	return &CreateProductHandler{log: log, cfg: cfg, repository: repository, kafkaProducer: kafkaProducer}
}

func (c *CreateProductHandler) Handle(ctx context.Context, command *CreateProduct) (*dtos.CreateProductResponseDto, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "CreateProductHandler.Handle")
	span.LogFields(log.String("ProductId", command.ProductID))
	defer span.Finish()

	product := CreateProductCommandToProductModel(command)

	product, err := c.repository.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	response := &dtos.CreateProductResponseDto{ProductID: product.ProductID}
	bytes, _ := json.Marshal(response)

	span.LogFields(log.String("CreateProductResponseDto", string(bytes)))

	return response, nil
}
