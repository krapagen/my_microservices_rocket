package config

type shipAssembledConsumerConfig struct {
	TopicName string `yaml:"topic" env:"SHIP_ASSEMBLED_TOPIC_NAME" env-default:"assembly.ship-assembled"`
	Group     string `yaml:"group_id" env:"SHIP_ASSEMBLED_GROUP_ID" env-default:"1"`
}

func (c *shipAssembledConsumerConfig) ShipAssembledTopicName() string {
	return c.TopicName
}

func (c *shipAssembledConsumerConfig) ShipAssembledGroup() string {
	return c.Group
}
