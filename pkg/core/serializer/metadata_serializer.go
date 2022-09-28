package serializer

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/metadata"
)

type MetadataSerializer interface {
	Serialize(meta metadata.Metadata) ([]byte, error)
	Deserialize(bytes []byte) (metadata.Metadata, error)
}
