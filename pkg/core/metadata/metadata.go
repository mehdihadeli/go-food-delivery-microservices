package metadata

import "time"

type Metadata map[string]interface{}

func (m Metadata) ExistsKey(key string) bool {
	_, exists := m[key]
	return exists
}

func (m Metadata) GetKey(key string) interface{} {
	val, exists := m[key]
	if !exists {
		return nil
	}

	return val
}

func (m Metadata) GetString(key string) string {
	val, ok := m.GetKey(key).(string)
	if ok {
		return val
	}

	return ""
}

func (m Metadata) GetTime(key string) time.Time {
	val, ok := m.GetKey(key).(time.Time)
	if ok {
		return val
	}

	return *new(time.Time)
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
