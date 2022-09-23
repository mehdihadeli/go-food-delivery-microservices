package typeMapper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypes(t *testing.T) {
	s := Test{A: 10}
	s2 := &Test{A: 10}

	q2 := TypeByName("*typeMapper.Test")
	q22 := InstanceByTypeName("*typeMapper.Test").(*Test)
	q222 := InstancePointerByTypeName("*typeMapper.Test").(*Test)
	q3 := TypeByNameAndImplementedInterface[ITest]("*typeMapper.Test")
	q4 := InstanceByTypeNameAndImplementedInterface[ITest]("*typeMapper.Test")

	c1 := GetTypeFromGeneric[Test]()
	c2 := GetTypeFromGeneric[*Test]()
	c3 := GetTypeFromGeneric[ITest]()

	d1 := GetType(Test{})
	d2 := GetType(&Test{})
	d3 := GetType((*ITest)(nil))

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

	assert.NotNil(t, d1)
	assert.NotNil(t, d2)
	assert.NotNil(t, d3)
	assert.NotNil(t, c1)
	assert.NotNil(t, c2)
	assert.NotNil(t, c3)
	assert.NotNil(t, q)
	assert.NotNil(t, q1)
	assert.NotNil(t, q11)
	assert.NotNil(t, q2)
	assert.NotNil(t, q22)
	assert.NotNil(t, q222)
	assert.NotNil(t, q3)
	assert.NotNil(t, q4)
	assert.NotNil(t, r)
	assert.NotNil(t, r2)
	assert.NotNil(t, r3)
	assert.NotNil(t, r4)
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

type ITest interface {
	Method1()
}

func (t *Test) Method1() {
}
