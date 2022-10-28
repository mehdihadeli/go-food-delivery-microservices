package gorm

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Gorm_Container(t *testing.T) {
	gorm, err := NewGormTestContainers().Start(context.Background(), t)
	require.NoError(t, err)

	assert.NotNil(t, gorm)
}
