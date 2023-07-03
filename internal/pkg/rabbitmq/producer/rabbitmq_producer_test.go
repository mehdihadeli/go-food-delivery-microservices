//go:build go1.18

package producer

import (
	"context"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/config/environemnt"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer/json"
	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/default_logger"
	types2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	config2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/types"
	testUtils "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/test/utils"
)

func Test_Publish_Message(t *testing.T) {
	testUtils.SkipCI(t)

	defaultLogger.SetupDefaultLogger()
	eventSerializer := serializer.NewDefaultEventSerializer(json.NewDefaultSerializer())

	ctx := context.Background()
	tp, err := tracing.NewOtelTracing(
		&config2.OpenTelemetryOptions{
			ServiceName:     "test",
			Enabled:         true,
			AlwaysOnSampler: true,
			JaegerExporterOptions: &config2.JaegerExporterOptions{
				AgentHost: "localhost",
				AgentPort: "6831",
			},
		},
		environemnt.Development,
	)
	if err != nil {
		return
	}
	defer tp.TracerProvider.Shutdown(ctx)

	conn, err := types.NewRabbitMQConnection(&config.RabbitmqOptions{
		RabbitmqHostOptions: &config.RabbitmqHostOptions{
			UserName: "guest",
			Password: "guest",
			HostName: "localhost",
			Port:     5672,
		},
	})
	require.NoError(t, err)

	rabbitmqProducer, err := NewRabbitMQProducer(
		conn,
		nil,
		defaultLogger.Logger,
		eventSerializer,
	)
	require.NoError(t, err)

	err = rabbitmqProducer.PublishMessage(ctx, NewProducerMessage("test"), nil)
	require.NoError(t, err)
}

type ProducerMessage struct {
	*types2.Message
	Data string
}

func NewProducerMessage(data string) *ProducerMessage {
	return &ProducerMessage{
		Data:    data,
		Message: types2.NewMessage(uuid.NewV4().String()),
	}
}
