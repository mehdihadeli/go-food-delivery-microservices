package consumers

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts/proto/kafka_messages"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	updatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products/commands/v1"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

type updateProductConsumer struct {
	*delivery.ProductConsumersBase
}

func NewUpdateProductConsumer(productConsumerBase *delivery.ProductConsumersBase) *updateProductConsumer {
	return &updateProductConsumer{productConsumerBase}
}

const (
	retryAttempts = 3
	retryDelay    = 300 * time.Millisecond
)

var (
	retryOptions = []retry.Option{retry.Attempts(retryAttempts), retry.Delay(retryDelay), retry.DelayType(retry.BackOffDelay)}
)

func (c *updateProductConsumer) Consume(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	c.Metrics.UpdateProductKafkaMessages.Inc()

	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "updateProductConsumer.Consume")
	defer span.Finish()

	msg := &kafka_messages.ProductUpdated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		c.Log.WarnMsg("proto.Unmarshal", err)
		c.CommitErrMessage(ctx, r, m)
		return
	}

	p := msg.GetProduct()

	productUUID, err := uuid.FromString(p.GetProductID())
	if err != nil {
		c.Log.WarnMsg("uuid.FromString", err)
		tracing.TraceErr(span, err)
		c.CommitErrMessage(ctx, r, m)
		return
	}

	command := updatingProductV1.NewUpdateProduct(productUUID, p.GetName(), p.GetDescription(), p.GetPrice())
	if err := c.Validator.StructCtx(ctx, command); err != nil {
		c.Log.WarnMsg("validate", err)
		tracing.TraceErr(span, err)
		c.CommitErrMessage(ctx, r, m)
		return
	}

	if err := retry.Do(func() error {
		_, err := mediatr.Send[*updatingProductV1.UpdateProductCommand, *mediatr.Unit](ctx, command)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		return nil
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		c.Log.WarnMsg("UpdateProductCommand.Handle", err)
		tracing.TraceErr(span, err)
		c.CommitErrMessage(ctx, r, m)
		return
	}

	c.CommitMessage(ctx, r, m)
}
