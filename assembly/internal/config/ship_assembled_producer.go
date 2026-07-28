package config

type shipAssembledProducerConfig struct {
	Topic string `yaml:"topic" env:"SHIP_ASSEMBLED_PRODUCER_TOPIC" env-default:"assembly.ship-assembled"`
}

func (c *shipAssembledProducerConfig) ShipAssembledProducerTopic() string {
	return c.Topic
}
