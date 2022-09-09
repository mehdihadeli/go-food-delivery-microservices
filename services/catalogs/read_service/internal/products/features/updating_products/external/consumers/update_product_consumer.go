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
	updatingProductV1 "github.com/mehdihadeli/store-golang-microservice-sample/services/catalogs/read_service/internal/products/features/updating_products/commands/v1"
	"github.com/opentracing/opentracing-go/log"
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
	span.LogFields(log.Object("Message", m))
	defer span.Finish()

	msg := &kafka_messages.ProductUpdated{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		unMarshalErr := customErrors.NewUnMarshalingErrorWrap(err, "[updateProductConsumer_Consume.Unmarshal] error in unMarshaling message")
		c.Log.Errorf(fmt.Sprintf("[updateProductConsumer_Consume.Unmarshal] err: %v", tracing.TraceWithErr(span, unMarshalErr)))
		c.CommitErrMessage(ctx, r, m)

		return
	}

	p := msg.GetProduct()

	productUUID, err := uuid.FromString(p.GetProductID())
	if err != nil {
		c.Log.WarnMsg("uuid.FromString", err)
		badRequestErr := customErrors.NewBadRequestErrorWrap(err, "[updateProductConsumer_Consume.uuid.FromString] error in the converting uuid")
		c.Log.Errorf(fmt.Sprintf("[updateProductConsumer_Consume.uuid.FromString] err: %v", tracing.TraceWithErr(span, badRequestErr)))
		c.CommitErrMessage(ctx, r, m)

		return
	}

	command := updatingProductV1.NewUpdateProduct(productUUID, p.GetName(), p.GetDescription(), p.GetPrice())
	if err := c.Validator.StructCtx(ctx, command); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(err, "[updateProductConsumer_Consume.StructCtx] command validation failed")
		c.Log.Errorf(fmt.Sprintf("[updateProductConsumer_Consume.StructCtx] err: {%v}", tracing.TraceWithErr(span, validationErr)))
		c.CommitErrMessage(ctx, r, m)

		return
	}

	if err := retry.Do(func() error {
		_, err := mediatr.Send[*updatingProductV1.UpdateProductCommand, *mediatr.Unit](ctx, command)
		return err
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		err = errors.WithMessage(err, "[updateProductConsumer_Consume.Send] error in sending UpdateProductCommand")
		c.Log.Errorw(fmt.Sprintf("[updateProductConsumer_Consume.Send] id: {%s}, err: {%v}", command.ProductID, tracing.TraceWithErr(span, err)), logger.Fields{"ProductId": command.ProductID})

		c.CommitErrMessage(ctx, r, m)
		return
	}

	c.CommitMessage(ctx, r, m)
}
