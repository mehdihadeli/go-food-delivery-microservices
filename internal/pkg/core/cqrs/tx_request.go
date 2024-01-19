package cqrs

// https://www.mohitkhare.com/blog/go-naming-conventions/

type TxRequest interface {
	Request

	isTxRequest() bool
}
