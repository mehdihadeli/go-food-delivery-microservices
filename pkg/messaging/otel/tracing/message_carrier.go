package tracing

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/metadata"

type MessageCarrier struct {
	meta *metadata.Metadata
}

func NewMessageCarrier(meta *metadata.Metadata) MessageCarrier {
	return MessageCarrier{meta: meta}
}

func (a MessageCarrier) Get(key string) string {
	return a.meta.GetString(key)
}

func (a MessageCarrier) Set(key string, value string) {
	a.meta.Set(key, value)
}

func (a MessageCarrier) Keys() []string {
	return a.meta.Keys()
}
