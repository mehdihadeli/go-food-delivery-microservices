package cqrs

type request struct{}

type Request interface {
	isRequest()
}

func NewRequest() Request {
	return &request{}
}

func (r *request) isRequest() {
}

func IsRequest(obj interface{}) bool {
	if _, ok := obj.(Request); ok {
		return true
	}

	return false
}
