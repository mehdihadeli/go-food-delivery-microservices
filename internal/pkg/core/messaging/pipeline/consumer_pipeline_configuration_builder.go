package pipeline

type ConsumerPipelineConfigurationBuilderFunc func(ConsumerPipelineConfigurationBuilder)

type ConsumerPipelineConfigurationBuilder interface {
	AddPipeline(pipeline ConsumerPipeline) ConsumerPipelineConfigurationBuilder
	Build() *ConsumerPipelineConfiguration
}

type consumerPipelineConfigurationBuilder struct {
	pipelineConfigurations *ConsumerPipelineConfiguration
}

func NewConsumerPipelineConfigurationBuilder() ConsumerPipelineConfigurationBuilder {
	return &consumerPipelineConfigurationBuilder{pipelineConfigurations: &ConsumerPipelineConfiguration{}}
}

func (c *consumerPipelineConfigurationBuilder) AddPipeline(
	pipeline ConsumerPipeline,
) ConsumerPipelineConfigurationBuilder {
	c.pipelineConfigurations.Pipelines = append(c.pipelineConfigurations.Pipelines, pipeline)
	return c
}

func (c *consumerPipelineConfigurationBuilder) Build() *ConsumerPipelineConfiguration {
	return c.pipelineConfigurations
}
