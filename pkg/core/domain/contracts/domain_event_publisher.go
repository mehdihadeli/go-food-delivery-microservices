package contracts

type IDomainEventPublisher interface {
	Publish(event *any)
	PublishAll(events ...*any)
}
