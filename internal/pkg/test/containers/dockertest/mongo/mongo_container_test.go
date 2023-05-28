package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Mongo_Container(t *testing.T) {
	mongo, err := NewMongoDockerTest().Start(context.Background(), t)
	require.NoError(t, err)

	assert.NotNil(t, mongo)
}
