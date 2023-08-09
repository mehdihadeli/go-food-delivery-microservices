package metadata

type Metadata map[string]interface{}

func (m Metadata) ExistsKey(key string) bool {
	_, exists := m[key]
	return exists
}

func (m Metadata) Get(key string) interface{} {
	val, exists := m[key]
	if !exists {
		return nil
	}

	return val
}

func (m Metadata) Set(key string, value interface{}) {
	m[key] = value
}

func (m Metadata) Keys() []string {
	i := 0
	r := make([]string, len(m))

	for k := range m {
		r[i] = k
		i++
	}

	return r
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
