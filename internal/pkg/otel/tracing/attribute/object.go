package attribute

import (
    "go.opentelemetry.io/otel/attribute"

    "github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/serializer/jsonSerializer"
)

// Object creates a KeyValue with a interface{} Value type.
func Object(k string, v interface{}) attribute.KeyValue {
    marshal, err := jsonSerializer.Marshal(&v)
    if err != nil {
        return attribute.KeyValue{}
    }
    return attribute.KeyValue{Key: attribute.Key(k), Value: attribute.StringValue(string(marshal))}
}
