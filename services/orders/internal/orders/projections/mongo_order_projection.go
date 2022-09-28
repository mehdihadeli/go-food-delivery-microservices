package projections

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/contracts/projection"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/es/models"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mapper"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/producer"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/contracts/repositories"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/dtos"
	creatingOrderEvents "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/events/domain/v1"
	v1 "github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/features/creating_order/events/integration/v1"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/orders/internal/orders/models/orders/read_models"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type mongoOrderProjection struct {
	mongoOrderRepository repositories.OrderReadRepository
	rabbitmqProducer     producer.Producer
	logger               logger.Logger
}

func NewMongoOrderProjection(mongoOrderRepository repositories.OrderReadRepository, rabbitmqProducer producer.Producer, logger logger.Logger) projection.IProjection {
	return &mongoOrderProjection{mongoOrderRepository: mongoOrderRepository, rabbitmqProducer: rabbitmqProducer, logger: logger}
}

func (m mongoOrderProjection) ProcessEvent(ctx context.Context, streamEvent *models.StreamEvent) error {
	// Handling and projecting event to elastic read model
	switch evt := streamEvent.Event.(type) {

	case *creatingOrderEvents.OrderCreatedV1:
		return m.onOrderCreated(ctx, evt)
	}

	return nil
}

func (m *mongoOrderProjection) onOrderCreated(ctx context.Context, evt *creatingOrderEvents.OrderCreatedV1) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "mongoOrderProjection.onOrderCreated")
	span.LogFields(log.String("OrderId", evt.OrderId.String()))
	span.LogFields(log.Object("Event", evt))
	defer span.Finish()

	items, err := mapper.Map[[]*read_models.ShopItemReadModel](evt.ShopItems)
	if err != nil {
		return errors.WrapIf(err, "[mongoOrderProjection_onOrderCreated.Map] error in mapping shopItems")
	}

	orderRead := read_models.NewOrderReadModel(evt.OrderId, items, evt.AccountEmail, evt.DeliveryAddress, evt.DeliveredTime)
	_, err = m.mongoOrderRepository.CreateOrder(ctx, orderRead)
	if err != nil {
		return errors.WrapIf(err, "[mongoOrderProjection_onOrderCreated.CreateOrder] error in creating order with mongoOrderRepository")
	}

	orderReadDto, err := mapper.Map[*dtos.OrderReadDto](orderRead)
	if err != nil {
		return tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[mongoOrderProjection_onOrderCreated.Map] error in mapping OrderReadDto"))
	}

	orderCreatedEvent := v1.NewOrderCreatedV1(orderReadDto)

	err = m.rabbitmqProducer.PublishMessage(ctx, orderCreatedEvent, nil)
	if err != nil {
		return tracing.TraceWithErr(span, customErrors.NewApplicationErrorWrap(err, "[mongoOrderProjection_onOrderCreated.PublishMessage] error in publishing OrderCreated integration event"))
	}

	m.logger.Infow(fmt.Sprintf("[mongoOrderProjection.onOrderCreated] OrderCreated message with messageId `%s` published to the rabbitmq broker", orderCreatedEvent.MessageId), logger.Fields{"MessageId": orderCreatedEvent.MessageId, "Id": orderCreatedEvent.OrderId})

	m.logger.Infow(fmt.Sprintf("[mongoOrderProjection.onOrderCreated] order with id '%s' created", orderCreatedEvent.Id), logger.Fields{"Id": orderRead.Id, "MessageId": orderCreatedEvent.MessageId})

	return nil
}
