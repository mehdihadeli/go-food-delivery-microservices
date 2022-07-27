package consumers

import (
	"context"
	"github.com/avast/retry-go"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/mediatr"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts/proto/kafka_messages"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	deletingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/deleting_products/commands/v1"
	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
	"time"
)

type deleteProductConsumer struct {
	*delivery.ProductConsumersBase
}

func NewDeleteProductConsumer(productConsumerBase *delivery.ProductConsumersBase) *deleteProductConsumer {
	return &deleteProductConsumer{productConsumerBase}
}

const (
	retryAttempts = 3
	retryDelay    = 300 * time.Millisecond
)

var (
	retryOptions = []retry.Option{retry.Attempts(retryAttempts), retry.Delay(retryDelay), retry.DelayType(retry.BackOffDelay)}
)

func (c *deleteProductConsumer) Consume(ctx context.Context, r *kafka.Reader, m kafka.Message) {

	c.Metrics.DeleteProductKafkaMessages.Inc()

	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "deleteProductConsumer.consume")
	defer span.Finish()

	msg := &kafka_messages.ProductDeleted{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		c.Log.WarnMsg("proto.Unmarshal", err)
		c.CommitErrMessage(ctx, r, m)
		return
	}

	productUUID, err := uuid.FromString(msg.GetProductID())
	if err != nil {
		c.Log.WarnMsg("uuid.FromString", err)
		tracing.TraceErr(span, err)
		c.CommitErrMessage(ctx, r, m)
		return
	}

	command := deletingProductV1.NewDeleteProduct(productUUID)
	if err := c.Validator.StructCtx(ctx, command); err != nil {
		c.Log.WarnMsg("validate", err)
		tracing.TraceErr(span, err)
		c.CommitErrMessage(ctx, r, m)
		return
	}

	if err := retry.Do(func() error {
		_, err := mediatr.Send[*mediatr.Unit, *deletingProductV1.DeleteProduct](ctx, command)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		return nil
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		c.Log.WarnMsg("DeleteProduct.Handle", err)
		c.CommitErrMessage(ctx, r, m)
		return
	}

	c.CommitMessage(ctx, r, m)
}
