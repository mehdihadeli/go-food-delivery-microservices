package kafka

// Config kafka config
type Config struct {
	Brokers    []string `mapstructure:"brokers"`
	GroupID    string   `mapstructure:"groupID" env:"GroupID"`
	InitTopics bool     `mapstructure:"initTopics" env:"InitTopics"`
}

// TopicConfig kafka topic config
type TopicConfig struct {
	TopicName         string `mapstructure:"topicName" env:"TopicName"`
	Partitions        int    `mapstructure:"partitions" env:"Partitions"`
	ReplicationFactor int    `mapstructure:"replicationFactor" env:"ReplicationFactor"`
}
