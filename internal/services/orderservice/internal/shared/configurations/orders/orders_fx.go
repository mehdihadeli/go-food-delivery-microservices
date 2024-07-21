package orders

import (
	"fmt"

	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/orders"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/shared/configurations/orders/infrastructure"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/services/orderservice/internal/shared/contracts"

	"go.opentelemetry.io/otel/metric"
	api "go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
)

// https://pmihaylov.com/shared-components-go-microservices/

var OrderServiceModule = fx.Module(
	"ordersfx",
	// Shared Modules
	config.Module,
	infrastructure.Module,

	// Features Modules
	orders.Module,

	// Other provides
	fx.Provide(configOrdersMetrics),
)

// ref: https://github.com/open-telemetry/opentelemetry-go/blob/main/example/prometheus/main.go

func configOrdersMetrics(
	cfg *config.Config,
	meter metric.Meter,
) (*contracts.OrdersMetrics, error) {
	if meter == nil {
		return nil, nil
	}

	appOptions := cfg.AppOptions
	successGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_success_grpc_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of success grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	errorGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_error_grpc_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of error grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	createOrderGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_create_order_grpc_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of create order grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	updateOrderGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_update_order_grpc_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of update order grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	payOrderGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_pay_order_grpc_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of pay order grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	submitOrderGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_submit_order_grpc_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of submit order grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	getOrderByIdGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_get_order_by_id_grpc_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of get order by id grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	getOrdersGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_get_orders_grpc_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of get orders grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	searchOrderGrpcRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_search_order_grpc_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of search order grpc requests"),
	)
	if err != nil {
		return nil, err
	}

	getOrdersHttpRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_get_orders_http_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of get orders http requests"),
	)
	if err != nil {
		return nil, err
	}

	createOrderHttpRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_create_order_http_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of create order http requests"),
	)
	if err != nil {
		return nil, err
	}

	updateOrderHttpRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_update_order_http_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of update order http requests"),
	)
	if err != nil {
		return nil, err
	}

	payOrderHttpRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_pay_order_http_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of pay order http requests"),
	)
	if err != nil {
		return nil, err
	}

	submitOrderHttpRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_submit_order_http_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of submit order http requests"),
	)
	if err != nil {
		return nil, err
	}

	getOrderByIdHttpRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_get_order_by_id_http_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of get order by id http requests"),
	)
	if err != nil {
		return nil, err
	}

	searchOrderHttpRequests, err := meter.Float64Counter(
		fmt.Sprintf("%s_search_order_http_requests_total", appOptions.ServiceName),
		api.WithDescription("The total number of search order http requests"),
	)
	if err != nil {
		return nil, err
	}

	deleteOrderRabbitMQMessages, err := meter.Float64Counter(
		fmt.Sprintf("%s_delete_order_rabbitmq_messages_total", appOptions.ServiceName),
		api.WithDescription("The total number of delete order rabbirmq messages"),
	)
	if err != nil {
		return nil, err
	}

	createOrderRabbitMQMessages, err := meter.Float64Counter(
		fmt.Sprintf("%s_create_order_rabbitmq_messages_total", appOptions.ServiceName),
		api.WithDescription("The total number of create order rabbirmq messages"),
	)
	if err != nil {
		return nil, err
	}

	updateOrderRabbitMQMessages, err := meter.Float64Counter(
		fmt.Sprintf("%s_update_order_rabbitmq_messages_total", appOptions.ServiceName),
		api.WithDescription("The total number of update order rabbirmq messages"),
	)
	if err != nil {
		return nil, err
	}

	return &contracts.OrdersMetrics{
		CreateOrderHttpRequests:     createOrderHttpRequests,
		SuccessGrpcRequests:         successGrpcRequests,
		ErrorGrpcRequests:           errorGrpcRequests,
		CreateOrderGrpcRequests:     createOrderGrpcRequests,
		UpdateOrderGrpcRequests:     updateOrderGrpcRequests,
		PayOrderGrpcRequests:        payOrderGrpcRequests,
		SubmitOrderGrpcRequests:     submitOrderGrpcRequests,
		GetOrderByIdGrpcRequests:    getOrderByIdGrpcRequests,
		GetOrdersGrpcRequests:       getOrdersGrpcRequests,
		SearchOrderGrpcRequests:     searchOrderGrpcRequests,
		GetOrdersHttpRequests:       getOrdersHttpRequests,
		UpdateOrderHttpRequests:     updateOrderHttpRequests,
		PayOrderHttpRequests:        payOrderHttpRequests,
		SubmitOrderHttpRequests:     submitOrderHttpRequests,
		GetOrderByIdHttpRequests:    getOrderByIdHttpRequests,
		SearchOrderHttpRequests:     searchOrderHttpRequests,
		DeleteOrderRabbitMQMessages: deleteOrderRabbitMQMessages,
		CreateOrderRabbitMQMessages: createOrderRabbitMQMessages,
		UpdateOrderRabbitMQMessages: updateOrderRabbitMQMessages,
	}, nil
}
