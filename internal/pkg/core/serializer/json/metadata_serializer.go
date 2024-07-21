package json

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/metadata"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/core/serializer"

	"emperror.dev/errors"
)

type DefaultMetadataJsonSerializer struct {
	serializer serializer.Serializer
}

func NewDefaultMetadataJsonSerializer(serializer serializer.Serializer) serializer.MetadataSerializer {
	return &DefaultMetadataJsonSerializer{serializer: serializer}
}

func (s *DefaultMetadataJsonSerializer) Serialize(meta metadata.Metadata) ([]byte, error) {
	if meta == nil {
		return nil, nil
	}

	marshal, err := s.serializer.Marshal(meta)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to marshal metadata")
	}

	return marshal, nil
}

func (s *DefaultMetadataJsonSerializer) Deserialize(bytes []byte) (metadata.Metadata, error) {
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
