package domain

import (
	"fmt"
	"github.com/goccy/go-reflect"
)

//TODO: Integrate with mediatr for publishing events

type defaultDomainEventPublisher struct{}

func NewDefaultDomainEventPublisher() *defaultDomainEventPublisher {
	return &defaultDomainEventPublisher{}
}

func (r *defaultDomainEventPublisher) Publish(event *any) {
	fmt.Println("event that is published :" + reflect.TypeOf(event).Name())
}

func (r *defaultDomainEventPublisher) PublishAll(events ...*any) {
	for _, event := range events {
		r.Publish(event)
	}
}
