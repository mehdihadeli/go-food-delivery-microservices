package messaging

import (
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/messaging/types"
	uuid "github.com/satori/go.uuid"
)

type FakeMessage struct {
	*types.Message
}

func NewFakeMessage() *FakeMessage {
	return &FakeMessage{Message: types.NewMessage(uuid.NewV4().String())}
}
