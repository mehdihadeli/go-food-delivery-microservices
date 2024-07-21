package postgresgorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Options_Name(t *testing.T) {
	assert.Equal(t, "gormOptions", optionName)
}
