package postgrespxg

import (
	"context"
	"testing"

	defaultLogger "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/logger/defaultlogger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Custom_PostgresPgx_Container(t *testing.T) {
	gorm, err := NewPostgresPgxContainers(
		defaultLogger.GetLogger(),
	).Start(context.Background(), t)
	require.NoError(t, err)

	assert.NotNil(t, gorm)
}
