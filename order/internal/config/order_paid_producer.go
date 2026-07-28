package config

type orderPaidProducerConfig struct {
	Topic string `yaml:"topic" env:"ORDER_PAID_PRODUCER_TOPIC" env-default:"order.paid"`
}

func (c *orderPaidProducerConfig) OrderPaidProducerTopic() string {
	return c.Topic
}
