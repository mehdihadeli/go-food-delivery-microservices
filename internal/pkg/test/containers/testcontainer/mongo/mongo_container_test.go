package mongo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Mongo_Container(t *testing.T) {
	mongo, err := NewMongoTestContainers().Start(context.Background(), t)
	require.NoError(t, err)

	assert.NotNil(t, mongo)
}
