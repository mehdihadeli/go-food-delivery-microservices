package typeMapper

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypes(t *testing.T) {
	s := Test{A: 10}
	fmt.Print(s)

	q := TypeByName("typeMapper.Test")
	w := TypeInstanceByName("typeMapper.Test").(Test)
	f := TypePointerInstanceByName("typeMapper.Test").(*Test)
	v := GenericInstanceTypeByName[Test]("typeMapper.Test")

	assert.NotNil(t, q)
	assert.NotNil(t, w)
	assert.NotNil(t, v)
	assert.NotNil(t, f)
}

type Test struct {
	A int
}
