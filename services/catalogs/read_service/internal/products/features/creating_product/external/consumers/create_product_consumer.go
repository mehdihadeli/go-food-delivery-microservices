package consumers

import (
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts/proto/kafka_messages"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	creatingProduct "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/commands/v1"
	"github.com/opentracing/opentracing-go/log"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

type createProductConsumer struct {
	*delivery.ProductConsumersBase
}

func NewCreateProductConsumer(productConsumerBase *delivery.ProductConsumersBase) *createProductConsumer {
	return &createProductConsumer{productConsumerBase}
}

const (
	retryAttempts = 3
	retryDelay    = 300 * time.Millisecond
)

var (
	retryOptions = []retry.Option{retry.Attempts(retryAttempts), retry.Delay(retryDelay), retry.DelayType(retry.BackOffDelay)}
)

func (c *createProductConsumer) Consume(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	c.Metrics.CreateProductKafkaMessages.Inc()
	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "createProductConsumer.Consume")
	span.LogFields(log.Object("Message", m))
	defer span.Finish()

	msg := &kafka_messages.ProductCreated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		unMarshalErr := customErrors.NewUnMarshalingErrorWrap(err, "[createProductConsumer_Consume.Unmarshal] error in unMarshaling message")
		c.Log.Errorf(fmt.Sprintf("[createProductConsumer_Consume.Unmarshal] err: %v", tracing.TraceWithErr(span, unMarshalErr)))
		c.CommitErrMessage(ctx, r, m)

		return
	}

	p := msg.GetProduct()

	command := v1.NewCreateProduct(p.GetProductID(), p.GetName(), p.GetDescription(), p.GetPrice(), p.GetCreatedAt().AsTime())
	if err := c.Validator.StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[createProductConsumer_Consume.StructCtx] command validation failed")
		c.Log.Errorf(fmt.Sprintf("[createProductConsumer_Consume.StructCtx] err: {%v}", tracing.TraceWithErr(span, validationErr)))
		c.CommitErrMessage(ctx, r, m)

		return
	}

	if err := retry.Do(func() error {
		_, err := mediatr.Send[*v1.CreateProduct, *creatingProduct.CreateProductResponseDto](ctx, command)
		return err
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		err = errors.WithMessage(err, "[createProductConsumer_Consume.Send] error in sending CreateProduct")
		c.Log.Errorw(fmt.Sprintf("[createProductConsumer_Consume.Send] id: {%s}, err: {%v}", command.ProductID, tracing.TraceWithErr(span, err)), logger.Fields{"ProductId": command.ProductID})
		c.CommitErrMessage(ctx, r, m)

		return
	}

	c.CommitMessage(ctx, r, m)
}
