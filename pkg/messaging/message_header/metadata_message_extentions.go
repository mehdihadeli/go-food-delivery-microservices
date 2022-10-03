package messageHeader

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/core/metadata"
	"time"
)

func GetCorrelationId(m metadata.Metadata) string {
	return m.GetString(CorrelationId)
}

func SetCorrelationId(m metadata.Metadata, val string) {
	m.Set(CorrelationId, val)
}

func GetMessageId(m metadata.Metadata) string {
	return m.GetString(MessageId)
}

func SetMessageId(m metadata.Metadata, val string) {
	m.Set(MessageId, val)
}

func GetMessageName(m metadata.Metadata) string {
	return m.GetString(Name)
}

func SetMessageName(m metadata.Metadata, val string) {
	m.Set(Name, val)
}

func GetMessageType(m metadata.Metadata) string {
	return m.GetString(Type)
}

func SetMessageType(m metadata.Metadata, val string) {
	m.Set(Type, val)
}

func SetMessageContentType(m metadata.Metadata, val string) {
	m.Set(Type, val)
}

func GetMessageContentType(m metadata.Metadata) string {
	return m.GetString(ContentType)
}

func GetMessageCreated(m metadata.Metadata) time.Time {
	return m.GetTime(Created)
}

func SetMessageCreated(m metadata.Metadata, val time.Time) {
	m.Set(Created, val)
}
