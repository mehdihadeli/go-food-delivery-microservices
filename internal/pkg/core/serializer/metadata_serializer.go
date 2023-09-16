package serializer

import (
	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/core/metadata"

	"emperror.dev/errors"
)

type MetadataSerializer interface {
	Serialize(meta metadata.Metadata) ([]byte, error)
	Deserialize(bytes []byte) (metadata.Metadata, error)
}

type DefaultMetadataSerializer struct {
	serializer Serializer
}

func NewDefaultMetadataSerializer(serializer Serializer) MetadataSerializer {
	return &DefaultMetadataSerializer{serializer: serializer}
}

func (s *DefaultMetadataSerializer) Serialize(meta metadata.Metadata) ([]byte, error) {
	if meta == nil {
		return nil, nil
	}

	marshal, err := s.serializer.Marshal(meta)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to marshal metadata")
	}

	return marshal, nil
}

func (s *DefaultMetadataSerializer) Deserialize(bytes []byte) (metadata.Metadata, error) {
	if bytes == nil {
		return nil, nil
	}

	var meta metadata.Metadata

	err := s.serializer.Unmarshal(bytes, &meta)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to unmarshal metadata")
	}

	return meta, nil
}
