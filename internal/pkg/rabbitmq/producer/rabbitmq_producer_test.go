//go:build.sh go1.18

package producer

import (
	"context"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/serializer/json"
	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/default_logger"
	types2 "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/messaging/types"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/otel/tracing"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/config"
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/rabbitmq/types"
)

func Test_Publish_Message(t *testing.T) {
	testUtils.SkipCI(t)
	ctx := context.Background()
	tp, err := tracing.AddOtelTracing(
		&otel.OpenTelemetryConfig{
			ServiceName:     "test",
			Enabled:         true,
			AlwaysOnSampler: true,
			JaegerExporterConfig: &otel.JaegerExporterConfig{
				AgentHost: "localhost",
				AgentPort: "6831",
			},
		},
	)
	if err != nil {
		return
	}
	defer tp.Shutdown(ctx)

	conn, err := types.NewRabbitMQConnection(ctx, &config.RabbitMQConfig{
		RabbitMqHostOptions: &config.RabbitMqHostOptions{
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
		json.NewJsonEventSerializer(),
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
