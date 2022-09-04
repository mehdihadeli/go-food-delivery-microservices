package consumers

import (
	"context"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/mehdihadeli/store-golang-microservice-sample/pkg/http/http_errors/custom_errors"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/logger"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/tracing"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/contracts/proto/kafka_messages"
	"github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/delivery"
	deletingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/deleting_products/commands/v1"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
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
	span.LogFields(log.Object("Message", m))
	defer span.Finish()

	msg := &kafka_messages.ProductDeleted{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		unMarshalErr := customErrors.NewUnMarshalingErrorWrap(err, "[deleteProductConsumer_Consume.Unmarshal] error in unMarshaling message")
		c.Log.Errorf(fmt.Sprintf("[deleteProductConsumer_Consume.Unmarshal] err: %v", tracing.TraceWithErr(span, unMarshalErr)))
		c.CommitErrMessage(ctx, r, m)

		return
	}

	productUUID, err := uuid.FromString(msg.GetProductID())
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[deleteProductConsumer_Consume.uuid.FromString] error in the converting uuid")
		c.Log.Errorf(fmt.Sprintf("[deleteProductConsumer_Consume.uuid.FromString] err: %v", tracing.TraceWithErr(span, badRequestErr)))

		c.CommitErrMessage(ctx, r, m)
		return
	}

	command := deletingProductV1.NewDeleteProductCommand(productUUID)
	if err := c.Validator.StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[deleteProductConsumer_Consume.StructCtx] command validation failed")
		c.Log.Errorf(fmt.Sprintf("[deleteProductConsumer_Consume.StructCtx] err: {%v}", tracing.TraceWithErr(span, validationErr)))

		c.CommitErrMessage(ctx, r, m)
		return
	}

	if err := retry.Do(func() error {
		_, err := mediatr.Send[*deletingProductV1.DeleteProductCommand, *mediatr.Unit](ctx, command)
		return err
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		err = errors.WithMessage(err, "[deleteProductConsumer_Consume.Send] error in sending DeleteProductCommand")
		c.Log.Errorw(fmt.Sprintf("[deleteProductConsumer_Consume.Send] id: {%s}, err: {%v}", command.ProductID, tracing.TraceWithErr(span, err)), logger.Fields{"ProductId": command.ProductID})

		c.CommitErrMessage(ctx, r, m)
		return
	}

	c.CommitMessage(ctx, r, m)
}
