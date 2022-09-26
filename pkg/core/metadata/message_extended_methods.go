package metadata

import (
	messageHeader "github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/message_header"
	"time"
)

func (m Metadata) GetCorrelationId() string {
	return m.GetString(messageHeader.CorrelationId)
}

func (m Metadata) SetCorrelationId(val string) {
	m.SetValue(messageHeader.CorrelationId, val)
}

func (m Metadata) GetMessageId() string {
	return m.GetString(messageHeader.MessageId)
}

func (m Metadata) SetMessageId(val string) {
	m.SetValue(messageHeader.MessageId, val)
}

func (m Metadata) GetMessageName() string {
	return m.GetString(messageHeader.Name)
}

func (m Metadata) SetMessageName(val string) {
	m.SetValue(messageHeader.Name, val)
}

func (m Metadata) GetMessageType() string {
	return m.GetString(messageHeader.Type)
}

func (m Metadata) SetMessageType(val string) {
	m.SetValue(messageHeader.Type, val)
}

func (m Metadata) GetMessageCreated() time.Time {
	return m.GetTime(messageHeader.Created)
}

func (m Metadata) SetMessageCreated(val time.Time) {
	m.SetValue(messageHeader.Created, val)
}
