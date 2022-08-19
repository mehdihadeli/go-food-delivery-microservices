package core

type Metadata map[string]interface{}

func MapToMetadata(data map[string]interface{}) *Metadata {
	m := Metadata(data)
	return &m
}

func MetadataToMap(meta *Metadata) map[string]interface{} {
	return *meta
}
