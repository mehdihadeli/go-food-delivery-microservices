//go:build unit
// +build unit

package json

import (
	"testing"

	"github.com/goccy/go-reflect"
	"github.com/stretchr/testify/assert"
)

type person struct {
	Name string
	Age  int
}

var currentSerializer = NewDefaultJsonSerializer()

func Test_Deserialize_Unstructured_Data_Into_Empty_Interface(t *testing.T) {
	// https://www.sohamkamani.com/golang/json/#decoding-json-to-maps---unstructured-data
	// https://developpaper.com/mapstructure-of-go/
	// https://pkg.go.dev/encoding/json#Unmarshal
	// when we assign an object type to interface is not pointer object or we don't assign interface, defaultLogger unmarshaler can't deserialize it to the object type and serialize it to map[string]interface{}

	// To unmarshal JSON into an interface value, Unmarshal stores map[string]interface{}
	var jsonMap interface{}

	marshal, err := currentSerializer.Marshal(person{"John", 30})
	if err != nil {
		return
	}

	err = currentSerializer.Unmarshal(marshal, &jsonMap)
	if err != nil {
		panic(err)
	}

	t.Log(jsonMap)

	for key, value := range jsonMap.(map[string]interface{}) {
		t.Log(key, value)
	}

	assert.True(t, reflect.TypeOf(jsonMap).Kind() == reflect.Map)
	assert.True(t, reflect.TypeOf(jsonMap) == reflect.TypeOf(map[string]interface{}(nil)))
	assert.True(t, jsonMap.(map[string]interface{})["ShortTypeName"] == "John")
	assert.True(t, jsonMap.(map[string]interface{})["Age"] == float64(30))
}

func Test_Deserialize_Unstructured_Data_Into_Map(t *testing.T) {
	// https://www.sohamkamani.com/golang/json/#decoding-json-to-maps---unstructured-data
	// https://developpaper.com/mapstructure-of-go/
	// https://pkg.go.dev/encoding/json#Unmarshal
	// when we assign an object type to interface is not pointer object or we don't assign interface, defaultLogger unmarshaler can't deserialize it to the object type and serialize it to map[string]interface{}

	// To unmarshal a JSON object into a map, Unmarshal first establishes a map to use. If the map is nil, Unmarshal allocates a new map. Otherwise Unmarshal reuses the existing map, keeping existing entries. Unmarshal then stores key-value pairs from the JSON object into the map.
	var jsonMap map[string]interface{}

	marshal, err := currentSerializer.Marshal(person{"John", 30})
	if err != nil {
		return
	}

	err = currentSerializer.Unmarshal(marshal, &jsonMap)
	if err != nil {
		panic(err)
	}

	t.Log(jsonMap)

	for key, value := range jsonMap {
		t.Log(key, value)
	}

	assert.True(t, reflect.TypeOf(jsonMap).Kind() == reflect.Map)
	assert.True(t, reflect.TypeOf(jsonMap) == reflect.TypeOf(map[string]interface{}(nil)))
	assert.True(t, jsonMap["ShortTypeName"] == "John")
	assert.True(t, jsonMap["Age"] == float64(30))
}

func Test_Deserialize_Structured_Data_Struct(t *testing.T) {
	// https://pkg.go.dev/encoding/json#Unmarshal
	// when we assign object to explicit struct type, defaultLogger unmarshaler can deserialize it to the struct

	// To unmarshal JSON into a struct, Unmarshal matches incoming object keys to the keys used by Marshal (either the struct field name or its tag), preferring an exact match but also accepting a case-insensitive match.
	var jsonMap person = person{}
	v := reflect.ValueOf(&jsonMap)
	if v.Elem().Kind() == reflect.Interface && v.NumMethod() == 0 {
		t.Log("deserialize to map[string]interface{}")
	} else {
		t.Log("deserialize to struct")
	}

	serializedObj := person{Name: "John", Age: 30}
	marshal, err := currentSerializer.Marshal(serializedObj)
	if err != nil {
		return
	}

	err = currentSerializer.Unmarshal(marshal, &jsonMap)
	if err != nil {
		panic(err)
	}

	assert.True(t, jsonMap.Name == "John")
	assert.True(t, jsonMap.Age == 30)
	assert.True(t, reflect.TypeOf(jsonMap) == reflect.TypeOf(person{}))
	assert.Equal(t, serializedObj, jsonMap)
}

func Test_Deserialize_Structured_Data_Struct2(t *testing.T) {
	// https://pkg.go.dev/encoding/json#Unmarshal
	// when we assign object to explicit struct type, defaultLogger unmarshaler can deserialize it to the struct

	// To unmarshal JSON into a struct, Unmarshal matches incoming object keys to the keys used by Marshal (either the struct field name or its tag), preferring an exact match but also accepting a case-insensitive match.
	var jsonMap interface{} = &person{}

	serializedObj := person{Name: "John", Age: 30}
	marshal, err := currentSerializer.Marshal(serializedObj)
	if err != nil {
		return
	}

	err = currentSerializer.Unmarshal(marshal, jsonMap)
	if err != nil {
		panic(err)
	}

	assert.True(t, jsonMap.(*person).Name == "John")
	assert.True(t, jsonMap.(*person).Age == 30)
	assert.True(t, reflect.TypeOf(jsonMap).Elem() == reflect.TypeOf(person{}))
}

func Test_Deserialize_Structured_Data_Pointer(t *testing.T) {
	// https://pkg.go.dev/encoding/json#Unmarshal
	// when we assign object to explicit struct type, defaultLogger unmarshaler can deserialize it to the struct

	// To unmarshal JSON into a pointer, Unmarshal first handles the case of the JSON being the JSON literal null. In that case, Unmarshal sets the pointer to nil. Otherwise, Unmarshal unmarshals the JSON into the value pointed at by the pointer. If the pointer is nil, Unmarshal allocates a new value for it to point to.To unmarshal JSON into a struct, Unmarshal matches incoming object keys to the keys used by Marshal (either the struct field name or its tag), preferring an exact match but also accepting a case-insensitive match.
	var jsonMap *person = &person{}
	// var jsonMap *person = nil

	serializedObj := person{Name: "John", Age: 30}
	marshal, err := currentSerializer.Marshal(serializedObj)
	if err != nil {
		return
	}

	err = currentSerializer.Unmarshal(marshal, jsonMap)
	if err != nil {
		panic(err)
	}

	assert.True(t, jsonMap.Name == "John")
	assert.True(t, jsonMap.Age == 30)
	assert.True(t, reflect.TypeOf(jsonMap).Elem() == reflect.TypeOf(person{}))
}

func Test_Decode_To_Map(t *testing.T) {
	var jsonMap map[string]interface{}

	serializedObj := person{Name: "John", Age: 30}
	marshal, err := currentSerializer.Marshal(serializedObj)
	if err != nil {
		return
	}

	// https://pkg.go.dev/encoding/json#Unmarshal
	// To unmarshal a JSON object into a map, Unmarshal first establishes a map to use. If the map is nil, Unmarshal allocates a new map. Otherwise Unmarshal reuses the existing map, keeping existing entries. Unmarshal then stores key-value pairs from the JSON object into the map.
	err = currentSerializer.UnmarshalToMap(marshal, &jsonMap)
	if err != nil {
		panic(err)
	}

	assert.True(t, jsonMap["ShortTypeName"] == "John")
	assert.True(t, jsonMap["Age"] == float64(30))
}
