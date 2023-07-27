package eventstoredb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EventStoreDB_Container(t *testing.T) {
	esdbInstance, err := NewEventstoreDBTestContainers().Start(context.Background(), t)
	require.NoError(t, err)

	assert.NotNil(t, esdbInstance)
}
