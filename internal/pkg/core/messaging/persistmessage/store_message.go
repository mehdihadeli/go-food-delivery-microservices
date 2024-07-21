package persistmessage

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type MessageDeliveryType int

const (
	Outbox   MessageDeliveryType = 1
	Inbox    MessageDeliveryType = 2
	Internal MessageDeliveryType = 4
)

type MessageStatus int

const (
	Stored    MessageStatus = 1
	Processed MessageStatus = 2
)

type StoreMessage struct {
	ID            uuid.UUID `gorm:"primaryKey"`
	DataType      string
	Data          string
	CreatedAt     time.Time `gorm:"default:current_timestamp"`
	RetryCount    int
	MessageStatus MessageStatus
	DeliveryType  MessageDeliveryType
}

func NewStoreMessage(
	id uuid.UUID,
	dataType string,
	data string,
	deliveryType MessageDeliveryType,
) *StoreMessage {
	return &StoreMessage{
		ID:            id,
		DataType:      dataType,
		Data:          data,
		CreatedAt:     time.Now(),
		MessageStatus: Stored,
		RetryCount:    0,
		DeliveryType:  deliveryType,
	}
}

func (sm *StoreMessage) ChangeState(messageStatus MessageStatus) {
	sm.MessageStatus = messageStatus
}

func (sm *StoreMessage) IncreaseRetry() {
	sm.RetryCount++
}

func (sm *StoreMessage) TableName() string {
	return "store_messages"
}
