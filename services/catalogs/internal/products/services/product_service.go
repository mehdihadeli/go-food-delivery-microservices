package service

import (
	kafkaClient "github.com/mehdihadeli/store-golang-microservice-sample/pkg/kafka"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/config"
	creating_product "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/deleting_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/getting_product_by_id"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/features/updating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/internal/products/infrastructure/repositories"
)

type ProductService struct {
	Commands *commands.ProductCommands
	Queries  *queries.ProductQueries
}

func NewProductService(log logger.Logger, cfg *config.Config, pgRepo repositories.ProductRepository, kafkaProducer kafkaClient.Producer) *ProductService {

	updateProductHandler := updating_product.NewUpdateProductHandler(log, cfg, pgRepo, kafkaProducer)
	createProductHandler := creating_product.NewCreateProductHandler(log, cfg, pgRepo, kafkaProducer)
	deleteProductHandler := deleting_product.NewDeleteProductHandler(log, cfg, pgRepo, kafkaProducer)

	getProductByIdHandler := getting_product_by_id.NewGetProductByIdHandler(log, cfg, pgRepo)

	productCommands := commands.NewProductCommands(createProductHandler, updateProductHandler, deleteProductHandler)
	productQueries := queries.NewProductQueries(getProductByIdHandler)

	return &ProductService{Commands: productCommands, Queries: productQueries}
}
