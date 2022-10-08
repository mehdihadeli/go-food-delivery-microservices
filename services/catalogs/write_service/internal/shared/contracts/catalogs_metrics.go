package contracts

import "go.opentelemetry.io/otel/metric/instrument/syncfloat64"

type CatalogsMetrics interface {
	CreateProductGrpcRequests() syncfloat64.Counter
	UpdateProductGrpcRequests() syncfloat64.Counter
	DeleteProductGrpcRequests() syncfloat64.Counter
	GetProductByIdGrpcRequests() syncfloat64.Counter
	SearchProductGrpcRequests() syncfloat64.Counter
	CreateProductHttpRequests() syncfloat64.Counter
	UpdateProductHttpRequests() syncfloat64.Counter
	DeleteProductHttpRequests() syncfloat64.Counter
	GetProductByIdHttpRequests() syncfloat64.Counter
	GetProductsHttpRequests() syncfloat64.Counter
	SearchProductHttpRequests() syncfloat64.Counter
	SuccessRabbitMQMessages() syncfloat64.Counter
	ErrorRabbitMQMessages() syncfloat64.Counter
	CreateProductRabbitMQMessages() syncfloat64.Counter
	UpdateProductRabbitMQMessages() syncfloat64.Counter
	DeleteProductRabbitMQMessages() syncfloat64.Counter
}
