package cqrs

type notification struct {
	TypeInfo
}

type Notification interface {
	isNotification()

	TypeInfo
}

func NewNotificationByT[T any]() Notification {
	n := &notification{
		TypeInfo: NewTypeInfoT[T](),
	}

	return n
}

func (c *notification) isNotification() {
}

func IsNotification(obj interface{}) bool {
	if _, ok := obj.(Notification); ok {
		return true
	}

	return false
}
