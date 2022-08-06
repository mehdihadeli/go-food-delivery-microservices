package consumers

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts/proto/kafka_messages"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	creatingProduct "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/creating_product/commands/v1"
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
	defer span.Finish()

	msg := &kafka_messages.ProductCreated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		c.Log.WarnMsg("proto.Unmarshal", err)
		tracing.TraceErr(span, err)
		c.CommitErrMessage(ctx, r, m)

		return
	}

	p := msg.GetProduct()
	command := v1.NewCreateProductCommand(p.GetProductID(), p.GetName(), p.GetDescription(), p.GetPrice(), p.GetCreatedAt().AsTime())
	if err := c.Validator.StructCtx(ctx, command); err != nil {
		tracing.TraceErr(span, err)
		c.Log.WarnMsg("validate", err)
		c.CommitErrMessage(ctx, r, m)
		return
	}

	if err := retry.Do(func() error {
		_, err := mediatr.Send[*v1.CreateProductCommand, *creatingProduct.CreateProductResponseDto](ctx, command)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		return nil
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		c.Log.WarnMsg("CreateProductCommand.Handle", err)
		tracing.TraceErr(span, err)
		c.CommitErrMessage(ctx, r, m)
		return
	}

	c.CommitMessage(ctx, r, m)
}
