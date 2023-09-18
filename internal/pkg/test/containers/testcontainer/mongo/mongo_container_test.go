package mongo

import (
	"context"
	"testing"

	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/default_logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Custom_Mongo_Container(t *testing.T) {
	defaultLogger.SetupDefaultLogger()

	mongo, err := NewMongoTestContainers(defaultLogger.Logger).Start(context.Background(), t)
	require.NoError(t, err)

	assert.NotNil(t, mongo)
}
