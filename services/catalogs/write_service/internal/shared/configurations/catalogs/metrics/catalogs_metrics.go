package metrics

import (
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/write_service/internal/shared/contracts"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
)

type catalogsMetrics struct {
	createProductGrpcRequests     syncfloat64.Counter
	updateProductGrpcRequests     syncfloat64.Counter
	deleteProductGrpcRequests     syncfloat64.Counter
	getProductByIdGrpcRequests    syncfloat64.Counter
	searchProductGrpcRequests     syncfloat64.Counter
	createProductHttpRequests     syncfloat64.Counter
	updateProductHttpRequests     syncfloat64.Counter
	deleteProductHttpRequests     syncfloat64.Counter
	getProductByIdHttpRequests    syncfloat64.Counter
	getProductsHttpRequests       syncfloat64.Counter
	searchProductHttpRequests     syncfloat64.Counter
	successRabbitMQMessages       syncfloat64.Counter
	errorRabbitMQMessages         syncfloat64.Counter
	createProductRabbitMQMessages syncfloat64.Counter
	updateProductRabbitMQMessages syncfloat64.Counter
	deleteProductRabbitMQMessages syncfloat64.Counter
}

func ConfigCatalogsMetrics(cfg *config.Config, meter metric.Meter) (contracts.CatalogsMetrics, error) {
	createProductGrpcRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_create_product_grpc_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of create product grpc requests"))
	if err != nil {
		return nil, err
	}

	updateProductGrpcRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_update_product_grpc_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of update product grpc requests"))
	if err != nil {
		return nil, err
	}

	deleteProductGrpcRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_delete_product_grpc_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of delete product grpc requests"))
	if err != nil {
		return nil, err
	}

	getProductByIdGrpcRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_get_product_by_id_grpc_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of get product by id grpc requests"))
	if err != nil {
		return nil, err
	}

	searchProductGrpcRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_search_product_grpc_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of search product grpc requests"))
	if err != nil {
		return nil, err
	}

	createProductRabbitMQMessages, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_create_product_rabbitmq_messages_total", cfg.ServiceName), instrument.WithDescription("The total number of create product rabbirmq messages"))
	if err != nil {
		return nil, err
	}

	updateProductRabbitMQMessages, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_update_product_rabbitmq_messages_total", cfg.ServiceName), instrument.WithDescription("The total number of update product rabbirmq messages"))
	if err != nil {
		return nil, err
	}

	deleteProductRabbitMQMessages, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_delete_product_rabbitmq_messages_total", cfg.ServiceName), instrument.WithDescription("The total number of delete product rabbirmq messages"))
	if err != nil {
		return nil, err
	}

	successRabbitMQMessages, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_search_product_rabbitmq_messages_total", cfg.ServiceName), instrument.WithDescription("The total number of success rabbitmq processed messages"))
	if err != nil {
		return nil, err
	}

	errorRabbitMQMessages, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_error_rabbitmq_processed_messages_total", cfg.ServiceName), instrument.WithDescription("The total number of error rabbitmq processed messages"))
	if err != nil {
		return nil, err
	}

	createProductHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_create_product_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of create product http requests"))
	if err != nil {
		return nil, err
	}

	updateProductHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_update_product_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of update product http requests"))
	if err != nil {
		return nil, err
	}

	deleteProductHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_delete_product_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of delete product http requests"))
	if err != nil {
		return nil, err
	}

	getProductByIdHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_get_product_by_id_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of get product by id http requests"))
	if err != nil {
		return nil, err
	}

	getProductsHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_get_products_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of get products http requests"))
	if err != nil {
		return nil, err
	}

	searchProductHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_search_product_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of search product http requests"))
	if err != nil {
		return nil, err
	}

	return &catalogsMetrics{
		createProductRabbitMQMessages: createProductRabbitMQMessages,
		getProductByIdGrpcRequests:    getProductByIdGrpcRequests,
		createProductGrpcRequests:     createProductGrpcRequests,
		createProductHttpRequests:     createProductHttpRequests,
		deleteProductRabbitMQMessages: deleteProductRabbitMQMessages,
		deleteProductGrpcRequests:     deleteProductGrpcRequests,
		deleteProductHttpRequests:     deleteProductHttpRequests,
		errorRabbitMQMessages:         errorRabbitMQMessages,
		getProductByIdHttpRequests:    getProductByIdHttpRequests,
		getProductsHttpRequests:       getProductsHttpRequests,
		searchProductGrpcRequests:     searchProductGrpcRequests,
		searchProductHttpRequests:     searchProductHttpRequests,
		successRabbitMQMessages:       successRabbitMQMessages,
		updateProductRabbitMQMessages: updateProductRabbitMQMessages,
		updateProductGrpcRequests:     updateProductGrpcRequests,
		updateProductHttpRequests:     updateProductHttpRequests,
	}, nil
}

func (c *catalogsMetrics) CreateProductGrpcRequests() syncfloat64.Counter {
	return c.createProductGrpcRequests
}

func (c *catalogsMetrics) UpdateProductGrpcRequests() syncfloat64.Counter {
	return c.updateProductGrpcRequests
}

func (c *catalogsMetrics) DeleteProductGrpcRequests() syncfloat64.Counter {
	return c.deleteProductGrpcRequests
}

func (c *catalogsMetrics) GetProductByIdGrpcRequests() syncfloat64.Counter {
	return c.getProductByIdGrpcRequests
}

func (c *catalogsMetrics) SearchProductGrpcRequests() syncfloat64.Counter {
	return c.searchProductGrpcRequests
}

func (c *catalogsMetrics) CreateProductHttpRequests() syncfloat64.Counter {
	return c.createProductHttpRequests
}

func (c *catalogsMetrics) UpdateProductHttpRequests() syncfloat64.Counter {
	return c.updateProductHttpRequests
}

func (c *catalogsMetrics) DeleteProductHttpRequests() syncfloat64.Counter {
	return c.deleteProductHttpRequests
}

func (c *catalogsMetrics) GetProductByIdHttpRequests() syncfloat64.Counter {
	return c.getProductByIdHttpRequests
}

func (c *catalogsMetrics) GetProductsHttpRequests() syncfloat64.Counter {
	return c.getProductsHttpRequests
}

func (c *catalogsMetrics) SearchProductHttpRequests() syncfloat64.Counter {
	return c.searchProductHttpRequests
}

func (c *catalogsMetrics) SuccessRabbitMQMessages() syncfloat64.Counter {
	return c.successRabbitMQMessages
}

func (c *catalogsMetrics) ErrorRabbitMQMessages() syncfloat64.Counter {
	return c.errorRabbitMQMessages
}

func (c *catalogsMetrics) CreateProductRabbitMQMessages() syncfloat64.Counter {
	return c.createProductRabbitMQMessages
}

func (c *catalogsMetrics) UpdateProductRabbitMQMessages() syncfloat64.Counter {
	return c.updateProductRabbitMQMessages
}

func (c *catalogsMetrics) DeleteProductRabbitMQMessages() syncfloat64.Counter {
	return c.deleteProductRabbitMQMessages
}
