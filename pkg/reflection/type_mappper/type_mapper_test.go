package typeMapper

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypes(t *testing.T) {
	s := Test{A: 10}
	s2 := &Test{A: 10}

	q2 := TypeByName("*typeMapper.Test")
	q22 := InstanceByTypeName("*typeMapper.Test").(*Test)
	q222 := InstancePointerByTypeName("*typeMapper.Test").(*Test)

	q := TypeByName("typeMapper.Test")
	q1 := InstanceByTypeName("typeMapper.Test").(Test)
	q11 := InstancePointerByTypeName("typeMapper.Test").(*Test)

	y := TypeByName("*Test")
	y1 := InstanceByTypeName("*Test").(*Test)
	y2 := InstancePointerByTypeName("*Test").(*Test)

	z := TypeByName("Test")
	z1 := InstanceByTypeName("Test").(Test)
	z2 := InstancePointerByTypeName("Test").(*Test)

	r := GenericInstanceByTypeName[*Test]("*typeMapper.Test")
	r2 := GenericInstanceByTypeName[Test]("typeMapper.Test")
	r3 := GenericInstanceByTypeName[Test]("Test")
	r4 := GenericInstanceByTypeName[*Test]("*Test")

	typeName := GetFullTypeName(s)
	typeName2 := GetFullTypeName(s2)
	typeName3 := GetTypeName(s2)
	typeName4 := GetTypeName(s)

	assert.Equal(t, "typeMapper.Test", typeName)
	assert.Equal(t, "*typeMapper.Test", typeName2)
	assert.Equal(t, "*Test", typeName3)
	assert.Equal(t, "Test", typeName4)

	q1.A = 100
	q22.A = 100

	fmt.Println(q.String())
	fmt.Println(q1)
	fmt.Println(q11)
	fmt.Println(q2.String())
	fmt.Println(q22)
	fmt.Println(q222)
	fmt.Println(q22.A)
	fmt.Println(q1.A)
	fmt.Println(r)
	fmt.Println(r2)
	fmt.Println(r3)
	fmt.Println(r4)
	assert.NotNil(t, q)
	assert.NotNil(t, q1)
	assert.NotNil(t, q11)
	assert.NotNil(t, q2)
	assert.NotNil(t, q22)
	assert.NotNil(t, q222)
	assert.NotNil(t, r)
	assert.NotNil(t, r2)
	assert.NotNil(t, y)
	assert.NotNil(t, y1)
	assert.NotNil(t, y2)
	assert.NotNil(t, z)
	assert.NotNil(t, z1)
	assert.NotNil(t, z2)
	assert.NotZero(t, q1.A)
	assert.NotZero(t, q22.A)
}

type Test struct {
	A int
}
