package shared

import (
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/contracts"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/deleting_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_product_by_id"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/getting_products"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/searching_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/products/features/updating_product"
	"github.com/pkg/errors"
)

func NewCatalogsMediator(log logger.Logger, cfg *config.Config, pgRepo contracts.ProductRepository, kafkaProducer kafkaClient.Producer) (*mediatr.Mediator, error) {

	md := mediatr.New()

	err := md.Register(
		creating_product.NewCreateProductHandler(log, cfg, pgRepo, kafkaProducer),
		updating_product.NewUpdateProductHandler(log, cfg, pgRepo, kafkaProducer),
		deleting_product.NewDeleteProductHandler(log, cfg, pgRepo, kafkaProducer),
		getting_product_by_id.NewGetProductByIdHandler(log, cfg, pgRepo),
		getting_products.NewGetProductsHandler(log, cfg, pgRepo),
		searching_product.NewSearchProductsHandler(log, cfg, pgRepo),
	)

	if err != nil {
		return nil, errors.Wrap(err, "error while registering handlers in the mediator")
	}

	return &md, nil
}
