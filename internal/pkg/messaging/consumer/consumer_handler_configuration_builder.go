package consumer

type ConsumerHandlerConfigurationBuilderFunc func(ConsumerHandlerConfigurationBuilder)

type ConsumerHandlerConfigurationBuilder interface {
	AddHandler(handler ConsumerHandler) ConsumerHandlerConfigurationBuilder
	Build() *ConsumerHandlersConfiguration
}

type consumerHandlerConfigurationBuilder struct {
	consumerHandlersConfiguration *ConsumerHandlersConfiguration
}

func NewConsumerHandlersConfigurationBuilder() ConsumerHandlerConfigurationBuilder {
	return &consumerHandlerConfigurationBuilder{consumerHandlersConfiguration: &ConsumerHandlersConfiguration{}}
}

func (c *consumerHandlerConfigurationBuilder) AddHandler(handler ConsumerHandler) ConsumerHandlerConfigurationBuilder {
	c.consumerHandlersConfiguration.Handlers = append(c.consumerHandlersConfiguration.Handlers, handler)
	return c
}

func (c *consumerHandlerConfigurationBuilder) Build() *ConsumerHandlersConfiguration {
	return c.consumerHandlersConfiguration
}
