package config

type kafkaConfig struct {
	Brokers []string `yaml:"brokers" env:"KAFKA_BROKERS" env-separator:"," env-default:"localhost:9092"`
}

func (c *kafkaConfig) KafkaBrokers() []string {
	return c.Brokers
}
