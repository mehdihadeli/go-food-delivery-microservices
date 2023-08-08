package gorm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Gorm_Container(t *testing.T) {
	gorm, err := NewGormDockerTest().Start(context.Background(), t)
	require.NoError(t, err)

	require.NotNil(t, gorm)
}
