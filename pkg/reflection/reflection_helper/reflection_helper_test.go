package reflectionHelper

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// ref: https://gist.github.com/drewolson/4771479
// https://stackoverflow.com/a/60598827/581476

type PersonPublic struct {
	Name string
	Age  int
}

type PersonPrivate struct {
	name string
	age  int
}

func (p *PersonPrivate) Name() string {
	return p.name
}

func (p *PersonPrivate) Age() int {
	return p.age
}

func Test_Field_Values_For_Exported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPublic{Name: "John", Age: 30}

	assert.Equal(t, "John", GetFieldValueByIndex(p, 0))
	assert.Equal(t, 30, GetFieldValueByIndex(p, 1))
}

func Test_Field_Values_For_Exported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPublic{Name: "John", Age: 30}

	assert.Equal(t, "John", GetFieldValueByIndex(p, 0))
	assert.Equal(t, 30, GetFieldValueByIndex(p, 1))
}

func Test_Field_Values_For_UnExported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPrivate{name: "John", age: 30}

	assert.Equal(t, "John", GetFieldValueByIndex(p, 0))
	assert.Equal(t, 30, GetFieldValueByIndex(p, 1))
}

func Test_Field_Values_For_UnExported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPrivate{name: "John", age: 30}

	assert.Equal(t, "John", GetFieldValueByIndex(p, 0))
	assert.Equal(t, 30, GetFieldValueByIndex(p, 1))
}

func Test_Set_Field_Value_For_Exported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPublic{}

	SetFieldValueByIndex(p, 0, "John")
	SetFieldValueByIndex(p, 1, 20)

	assert.Equal(t, "John", p.Name)
	assert.Equal(t, 20, p.Age)
}

func Test_Set_Field_Value_For_Exported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPublic{}

	SetFieldValueByIndex(&p, 0, "John")
	SetFieldValueByIndex(&p, 1, 20)

	assert.Equal(t, "John", p.Name)
	assert.Equal(t, 20, p.Age)
}

func Test_Set_Field_Value_For_UnExported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPrivate{}

	SetFieldValueByIndex(p, 0, "John")
	SetFieldValueByIndex(p, 1, 20)

	assert.Equal(t, "John", p.name)
	assert.Equal(t, 20, p.age)
}

func Test_Set_Field_Value_For_UnExported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPrivate{}

	SetFieldValueByIndex(&p, 0, "John")
	SetFieldValueByIndex(&p, 1, 20)

	assert.Equal(t, "John", p.name)
	assert.Equal(t, 20, p.age)
}

func Test_Get_Field_Value_For_Exported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPublic{Name: "John", Age: 20}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(p).Elem()
	name := GetFieldValue(v.FieldByName("Name")).Interface()
	age := GetFieldValue(v.FieldByName("Age")).Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 20, age)
}

func Test_Get_Field_Value_For_UnExported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPrivate{name: "John", age: 30}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(p).Elem()
	name := GetFieldValue(v.FieldByName("name")).Interface()
	age := GetFieldValue(v.FieldByName("age")).Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 30, age)
}

func Test_Get_Field_Value_For_Exported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPublic{Name: "John", Age: 20}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(&p).Elem()
	name := GetFieldValue(v.FieldByName("Name")).Interface()
	age := GetFieldValue(v.FieldByName("Age")).Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 20, age)
}

func Test_Get_Field_Value_For_UnExported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPrivate{name: "John", age: 20}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(&p).Elem()
	name := GetFieldValue(v.FieldByName("name")).Interface()
	age := GetFieldValue(v.FieldByName("age")).Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 20, age)
}

func Test_Set_Field_For_Exported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPublic{}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(p).Elem()
	name := GetFieldValue(v.FieldByName("Name"))
	age := GetFieldValue(v.FieldByName("Age"))

	SetFieldValue(name, "John")
	SetFieldValue(age, 20)

	assert.Equal(t, "John", name.Interface())
	assert.Equal(t, 20, age.Interface())
}

func Test_Set_Field_For_UnExported_Fields_And_Addressable_Struct(t *testing.T) {
	p := &PersonPrivate{}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(p).Elem()
	name := GetFieldValue(v.FieldByName("name"))
	age := GetFieldValue(v.FieldByName("age"))

	SetFieldValue(name, "John")
	SetFieldValue(age, 20)

	assert.Equal(t, "John", name.Interface())
	assert.Equal(t, 20, age.Interface())
}

func Test_Set_Field_For_Exported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPublic{}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(&p).Elem()
	name := GetFieldValue(v.FieldByName("Name"))
	age := GetFieldValue(v.FieldByName("Age"))

	SetFieldValue(name, "John")
	SetFieldValue(age, 20)

	assert.Equal(t, "John", name.Interface())
	assert.Equal(t, 20, age.Interface())
}

func Test_Set_Field_For_UnExported_Fields_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPrivate{}

	//field by name only work on struct not pointer type so we should get Elem()
	v := reflect.ValueOf(&p).Elem()
	name := GetFieldValue(v.FieldByName("name"))
	age := GetFieldValue(v.FieldByName("age"))

	SetFieldValue(name, "John")
	SetFieldValue(age, 20)

	assert.Equal(t, "John", name.Interface())
	assert.Equal(t, 20, age.Interface())
}

func Test_Get_Unexported_Field_From_Method_And_Addressable_Struct(t *testing.T) {
	p := &PersonPrivate{name: "John", age: 20}
	name := GetFieldValueFromMethodAndObject(p, "Name")

	assert.Equal(t, "John", name.Interface())
}

func Test_Get_Unexported_Field_From_Method_And_UnAddressable_Struct(t *testing.T) {
	p := PersonPrivate{name: "John", age: 20}
	name := GetFieldValueFromMethodAndObject(p, "Name")

	assert.Equal(t, "John", name.Interface())
}

func Test_Convert_NoPointer_Type_To_Pointer_Type_With_Addr(t *testing.T) {
	//https://www.geeksforgeeks.org/reflect-addr-function-in-golang-with-examples/

	p := PersonPrivate{name: "John", age: 20}
	v := reflect.ValueOf(&p).Elem()
	pointerType := v.Addr()
	name := pointerType.MethodByName("Name").Call(nil)[0].Interface()
	age := pointerType.MethodByName("Age").Call(nil)[0].Interface()

	assert.Equal(t, "John", name)
	assert.Equal(t, 20, age)
}
