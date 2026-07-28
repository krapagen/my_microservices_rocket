package config

type orderPaidConsumerConfig struct {
	TopicName string `yaml:"topic" env:"ORDER_PAID_TOPIC_NAME" env-default:"order.paid"`
	Group     string `yaml:"group_id" env:"ORDER_PAID_GROUP_ID" env-default:"1"`
}

func (c *orderPaidConsumerConfig) OrderPaidTopicName() string {
	return c.TopicName
}

func (c *orderPaidConsumerConfig) OrderPaidGroup() string {
	return c.Group
}
