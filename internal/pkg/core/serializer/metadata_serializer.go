package serializer

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/metadata"
)

type MetadataSerializer interface {
	Serialize(meta metadata.Metadata) ([]byte, error)
	Deserialize(bytes []byte) (metadata.Metadata, error)
}
