package metrics

import (
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/config"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/shared/contracts"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
)

type ordersServiceMetrics struct {
	successGrpcRequests syncfloat64.Counter
	errorGrpcRequests   syncfloat64.Counter

	createOrderGrpcRequests  syncfloat64.Counter
	updateOrderGrpcRequests  syncfloat64.Counter
	payOrderGrpcRequests     syncfloat64.Counter
	submitOrderGrpcRequests  syncfloat64.Counter
	getOrderByIdGrpcRequests syncfloat64.Counter
	getOrdersGrpcRequests    syncfloat64.Counter
	searchOrderGrpcRequests  syncfloat64.Counter

	successHttpRequests syncfloat64.Counter
	errorHttpRequests   syncfloat64.Counter

	createOrderHttpRequests  syncfloat64.Counter
	updateOrderHttpRequests  syncfloat64.Counter
	payOrderHttpRequests     syncfloat64.Counter
	submitOrderHttpRequests  syncfloat64.Counter
	getOrderByIdHttpRequests syncfloat64.Counter
	searchOrderHttpRequests  syncfloat64.Counter
	getOrdersHttpRequests    syncfloat64.Counter

	successRabbitMQMessages syncfloat64.Counter
	errorRabbitMQMessages   syncfloat64.Counter

	createOrderRabbitMQMessages syncfloat64.Counter
	updateOrderRabbitMQMessages syncfloat64.Counter
	deleteOrderRabbitMQMessages syncfloat64.Counter
}

func ConfigOrdersMetrics(cfg *config.Config, meter metric.Meter) (contracts.OrdersMetrics, error) {
	successGrpcRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_success_grpc_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of success grpc requests"))
	if err != nil {
		return nil, err
	}

	errorGrpcRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_error_grpc_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of error grpc requests"))
	if err != nil {
		return nil, err
	}

	createOrderGrpcRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_create_order_grpc_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of create order grpc requests"))
	if err != nil {
		return nil, err
	}

	updateOrderGrpcRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_update_order_grpc_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of update order grpc requests"))
	if err != nil {
		return nil, err
	}

	payOrderGrpcRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_pay_order_grpc_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of pay order grpc requests"))
	if err != nil {
		return nil, err
	}

	submitOrderGrpcRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_submit_order_grpc_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of submit order grpc requests"))
	if err != nil {
		return nil, err
	}

	getOrderByIdGrpcRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_get_order_by_id_grpc_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of get order by id grpc requests"))
	if err != nil {
		return nil, err
	}

	getOrdersGrpcRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_get_orders_grpc_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of get orders grpc requests"))
	if err != nil {
		return nil, err
	}

	searchOrderGrpcRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_search_order_grpc_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of search order grpc requests"))
	if err != nil {
		return nil, err
	}

	getOrdersHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_get_orders_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of get orders http requests"))
	if err != nil {
		return nil, err
	}

	successHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_success_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of success http requests"))
	if err != nil {
		return nil, err
	}

	errorHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_error_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of error http requests"))
	if err != nil {
		return nil, err
	}

	createOrderHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_create_order_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of create order http requests"))
	if err != nil {
		return nil, err
	}

	updateOrderHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_update_order_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of update order http requests"))
	if err != nil {
		return nil, err
	}

	payOrderHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_pay_order_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of pay order http requests"))
	if err != nil {
		return nil, err
	}

	submitOrderHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_submit_order_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of submit order http requests"))
	if err != nil {
		return nil, err
	}

	getOrderByIdHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_get_order_by_id_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of get order by id http requests"))
	if err != nil {
		return nil, err
	}

	searchOrderHttpRequests, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_search_order_http_requests_total", cfg.ServiceName), instrument.WithDescription("The total number of search order http requests"))
	if err != nil {
		return nil, err
	}

	deleteOrderRabbitMQMessages, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_delete_order_rabbitmq_messages_total", cfg.ServiceName), instrument.WithDescription("The total number of delete order rabbirmq messages"))
	if err != nil {
		return nil, err
	}

	createOrderRabbitMQMessages, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_create_order_rabbitmq_messages_total", cfg.ServiceName), instrument.WithDescription("The total number of create order rabbirmq messages"))
	if err != nil {
		return nil, err
	}

	updateOrderRabbitMQMessages, err := meter.SyncFloat64().Counter(fmt.Sprintf("%s_update_order_rabbitmq_messages_total", cfg.ServiceName), instrument.WithDescription("The total number of update order rabbirmq messages"))
	if err != nil {
		return nil, err
	}

	return &ordersServiceMetrics{
		createOrderHttpRequests:     createOrderHttpRequests,
		successGrpcRequests:         successGrpcRequests,
		errorGrpcRequests:           errorGrpcRequests,
		createOrderGrpcRequests:     createOrderGrpcRequests,
		updateOrderGrpcRequests:     updateOrderGrpcRequests,
		payOrderGrpcRequests:        payOrderGrpcRequests,
		submitOrderGrpcRequests:     submitOrderGrpcRequests,
		getOrderByIdGrpcRequests:    getOrderByIdGrpcRequests,
		getOrdersGrpcRequests:       getOrdersGrpcRequests,
		searchOrderGrpcRequests:     searchOrderGrpcRequests,
		getOrdersHttpRequests:       getOrdersHttpRequests,
		successHttpRequests:         successHttpRequests,
		errorHttpRequests:           errorHttpRequests,
		updateOrderHttpRequests:     updateOrderHttpRequests,
		payOrderHttpRequests:        payOrderHttpRequests,
		submitOrderHttpRequests:     submitOrderHttpRequests,
		getOrderByIdHttpRequests:    getOrderByIdHttpRequests,
		searchOrderHttpRequests:     searchOrderHttpRequests,
		deleteOrderRabbitMQMessages: deleteOrderRabbitMQMessages,
		createOrderRabbitMQMessages: createOrderRabbitMQMessages,
		updateOrderRabbitMQMessages: updateOrderRabbitMQMessages,
	}, nil
}

func (o ordersServiceMetrics) SuccessGrpcRequests() syncfloat64.Counter {
	return o.searchOrderGrpcRequests
}

func (o ordersServiceMetrics) ErrorGrpcRequests() syncfloat64.Counter {
	return o.errorGrpcRequests
}

func (o ordersServiceMetrics) CreateOrderGrpcRequests() syncfloat64.Counter {
	return o.createOrderGrpcRequests
}

func (o ordersServiceMetrics) UpdateOrderGrpcRequests() syncfloat64.Counter {
	return o.updateOrderGrpcRequests
}

func (o ordersServiceMetrics) PayOrderGrpcRequests() syncfloat64.Counter {
	return o.updateOrderGrpcRequests
}

func (o ordersServiceMetrics) SubmitOrderGrpcRequests() syncfloat64.Counter {
	return o.submitOrderGrpcRequests
}

func (o ordersServiceMetrics) GetOrderByIdGrpcRequests() syncfloat64.Counter {
	return o.getOrderByIdGrpcRequests
}

func (o ordersServiceMetrics) GetOrdersGrpcRequests() syncfloat64.Counter {
	return o.getOrdersGrpcRequests
}

func (o ordersServiceMetrics) SearchOrderGrpcRequests() syncfloat64.Counter {
	return o.searchOrderGrpcRequests
}

func (o ordersServiceMetrics) SuccessHttpRequests() syncfloat64.Counter {
	return o.successHttpRequests
}

func (o ordersServiceMetrics) ErrorHttpRequests() syncfloat64.Counter {
	return o.errorHttpRequests
}

func (o ordersServiceMetrics) CreateOrderHttpRequests() syncfloat64.Counter {
	return o.createOrderHttpRequests
}

func (o ordersServiceMetrics) UpdateOrderHttpRequests() syncfloat64.Counter {
	return o.updateOrderHttpRequests
}

func (o ordersServiceMetrics) PayOrderHttpRequests() syncfloat64.Counter {
	return o.payOrderHttpRequests
}

func (o ordersServiceMetrics) SubmitOrderHttpRequests() syncfloat64.Counter {
	return o.submitOrderHttpRequests
}

func (o ordersServiceMetrics) GetOrderByIdHttpRequests() syncfloat64.Counter {
	return o.getOrderByIdHttpRequests
}

func (o ordersServiceMetrics) SearchOrderHttpRequests() syncfloat64.Counter {
	return o.searchOrderHttpRequests
}

func (o ordersServiceMetrics) GetOrdersHttpRequests() syncfloat64.Counter {
	return o.getOrdersHttpRequests
}

func (o ordersServiceMetrics) SuccessRabbitMQMessages() syncfloat64.Counter {
	return o.successRabbitMQMessages
}

func (o ordersServiceMetrics) ErrorRabbitMQMessages() syncfloat64.Counter {
	return o.errorRabbitMQMessages
}

func (o ordersServiceMetrics) CreateOrderRabbitMQMessages() syncfloat64.Counter {
	return o.createOrderRabbitMQMessages
}

func (o ordersServiceMetrics) UpdateOrderRabbitMQMessages() syncfloat64.Counter {
	return o.updateOrderRabbitMQMessages
}

func (o ordersServiceMetrics) DeleteOrderRabbitMQMessages() syncfloat64.Counter {
	return o.deleteOrderRabbitMQMessages
}
