package shared

import (
	"fmt"
	"github.com/adzeitor/mediatr"
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/deleting_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/getting_product_by_id"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/updating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/infrastructure/repositories"
)

func NewMediator(log logger.Logger, cfg *config.Config, pgRepo repositories.ProductRepository, kafkaProducer kafkaClient.Producer) (*mediatr.Mediator, error) {

	md := mediatr.New()

	err := md.Register(creating_product.NewCreateProductHandler(log, cfg, pgRepo, kafkaProducer))
	if err != nil {
		return nil, fmt.Errorf("error in registering 'CreateProductHandler': %w", err)
	}
	err2 := md.Register(updating_product.NewUpdateProductHandler(log, cfg, pgRepo, kafkaProducer))
	if err2 != nil {
		return nil, fmt.Errorf("error in registering 'UpdateProductHandler': %w", err2)
	}
	err3 := md.Register(deleting_product.NewDeleteProductHandler(log, cfg, pgRepo, kafkaProducer))
	if err3 != nil {
		return nil, fmt.Errorf("error in registering 'DeleteProductHandler': %w", err3)
	}
	err4 := md.Register(getting_product_by_id.NewGetProductByIdHandler(log, cfg, pgRepo))
	if err4 != nil {
		return nil, fmt.Errorf("error in registering 'GetProductByIdHandler': %w", err4)
	}

	if err != nil {
		return nil, fmt.Errorf("create mediator: %w", err)
	}
	return &md, nil
}
