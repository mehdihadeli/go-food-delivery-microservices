package catalogs

import (
	"fmt"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/products/data/uow"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/configurations/catalogs/infrastructure"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/services/catalogwriteservice/internal/shared/contracts"

	"go.opentelemetry.io/otel/metric"
	api "go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
)

// https://pmihaylov.com/shared-components-go-microservices/
var CatalogsServiceModule = fx.Module(
	"catalogsfx",
	// Shared Modules
	config.Module,
	infrastructure.Module,

	// Features Modules
	products.Module,

	// Other provides
	fx.Provide(provideCatalogsMetrics),
	fx.Provide(uow.NewCatalogsUnitOfWork),
)

// ref: https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go
func provideCatalogsMetrics(
	cfg *config.AppOptions,
	meter metric.Meter,
) (*contracts.CatalogsMetrics, error) {
	if meter == nil {
		return nil, nil
	}

	createProductGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_create_product_grpc_requests_total", cfg.ServiceName),
		api.WithDescription("The total number of create product grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	updateProductGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_update_product_grpc_requests_total", cfg.ServiceName),
		api.WithDescription("The total number of update product grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	deleteProductGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_delete_product_grpc_requests_total", cfg.ServiceName),
		api.WithDescription("The total number of delete product grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	getProductByIdGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_get_product_by_id_grpc_requests_total", cfg.ServiceName),
		api.WithDescription("The total number of get product by id grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	searchProductGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_search_product_grpc_requests_total", cfg.ServiceName),
		api.WithDescription("The total number of search product grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	createProductRabbitMQMessages, err := meter.Float64Counter(
		fmt.Sprintf("%s_create_product_rabbitmq_messages_total", cfg.ServiceName),
		api.WithDescription("The total number of create product rabbirmq messages"),
	)
	if err != nil {
		return nil, err
	}

	updateProductRabbitMQMessages, err := meter.Float64Counter(
		fmt.Sprintf("%s_update_product_rabbitmq_messages_total", cfg.ServiceName),
		api.WithDescription("The total number of update product rabbirmq messages"),
	)
	if err != nil {
		return nil, err
	}

	deleteProductRabbitMQMessages, err := meter.Float64Counter(
		fmt.Sprintf("%s_delete_product_rabbitmq_messages_total", cfg.ServiceName),
		api.WithDescription("The total number of delete product rabbirmq messages"),
	)
	if err != nil {
		return nil, err
	}

	successRabbitMQMessages, err := meter.Float64Counter(
		fmt.Sprintf("%s_search_product_rabbitmq_messages_total", cfg.ServiceName),
		api.WithDescription("The total number of success rabbitmq processed messages"),
	)
	if err != nil {
		return nil, err
	}

	errorRabbitMQMessages, err := meter.Float64Counter(
		fmt.Sprintf("%s_error_rabbitmq_processed_messages_total", cfg.ServiceName),
		api.WithDescription("The total number of error rabbitmq processed messages"),
	)
	if err != nil {
		return nil, err
	}

	createProductHttpRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_create_product_http_requests_total", cfg.ServiceName),
		api.WithDescription("The total number of create product http requests"),
	)
	if err != nil {
		return nil, err
	}

	updateProductHttpRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_update_product_http_requests_total", cfg.ServiceName),
		api.WithDescription("The total number of update product http requests"),
	)
	if err != nil {
		return nil, err
	}

	deleteProductHttpRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_delete_product_http_requests_total", cfg.ServiceName),
		api.WithDescription("The total number of delete product http requests"),
	)
	if err != nil {
		return nil, err
	}

	getProductByIdHttpRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_get_product_by_id_http_requests_total", cfg.ServiceName),
		api.WithDescription("The total number of get product by id http requests"),
	)
	if err != nil {
		return nil, err
	}

	getProductsHttpRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_get_products_http_requests_total", cfg.ServiceName),
		api.WithDescription("The total number of get products http requests"),
	)
	if err != nil {
		return nil, err
	}

	searchProductHttpRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_search_product_http_requests_total", cfg.ServiceName),
		api.WithDescription("The total number of search product http requests"),
	)
	if err != nil {
		return nil, err
	}

	return &contracts.CatalogsMetrics{
		CreateProductRabbitMQMessages: createProductRabbitMQMessages,
		GetProductByIdGrpcRequests:    getProductByIdGrpcRequests,
		CreateProductGrpcRequests:     createProductGrpcRequests,
		CreateProductHttpRequests:     createProductHttpRequests,
		DeleteProductRabbitMQMessages: deleteProductRabbitMQMessages,
		DeleteProductGrpcRequests:     deleteProductGrpcRequests,
		DeleteProductHttpRequests:     deleteProductHttpRequests,
		ErrorRabbitMQMessages:         errorRabbitMQMessages,
		GetProductByIdHttpRequests:    getProductByIdHttpRequests,
		GetProductsHttpRequests:       getProductsHttpRequests,
		SearchProductGrpcRequests:     searchProductGrpcRequests,
		SearchProductHttpRequests:     searchProductHttpRequests,
		SuccessRabbitMQMessages:       successRabbitMQMessages,
		UpdateProductRabbitMQMessages: updateProductRabbitMQMessages,
		UpdateProductGrpcRequests:     updateProductGrpcRequests,
		UpdateProductHttpRequests:     updateProductHttpRequests,
	}, nil
}
