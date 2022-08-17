package esSerializer

import "github.com/mehdihadeli/store-golang-microservice-sample/pkg/core"

type MetadataSerializer interface {
	Serialize(meta *core.Metadata) ([]byte, error)
	Deserialize(bytes []byte) (*core.Metadata, error)
}
