package cqrs

// https://www.mohitkhare.com/blog/go-naming-conventions/
type ITxRequest interface {
	IsTxRequest() bool
}
