package cqrs

// https://www.mohitkhare.com/blog/go-naming-conventions/
// https://github.com/EventStore/EventStore-Client-Go/blob/master/esdb/position.go
type TxRequest interface {
	Request

	isTxRequest()
}
