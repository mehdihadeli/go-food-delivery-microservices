package core

import "emperror.dev/errors"

type Metadata map[string]interface{}

func (m Metadata) ExistsKey(key string) bool {
	_, exists := m[key]
	return exists
}

func (m Metadata) GetKey(key string) (interface{}, error) {
	val, exists := m[key]
	if !exists {
		return nil, errors.New("key not found")
	}

	return val, nil
}

func (m Metadata) SetValue(key string, value interface{}) {
	m[key] = value
}

func MapToMetadata(data map[string]interface{}) Metadata {
	m := Metadata(data)
	return m
}

func MetadataToMap(meta Metadata) map[string]interface{} {
	return meta
}

func FromMetadata(m Metadata) Metadata {

	if m == nil {
		return Metadata{}
	}
	return m
}
