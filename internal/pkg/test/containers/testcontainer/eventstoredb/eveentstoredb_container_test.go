package eventstoredb

import (
	"context"
	"testing"

	defaultLogger "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/logger/default_logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Custom_EventStoreDB_Container(t *testing.T) {
	esdbInstance, err := NewEventstoreDBTestContainers(defaultLogger.GetLogger()).Start(context.Background(), t)
	require.NoError(t, err)

	assert.NotNil(t, esdbInstance)
}
