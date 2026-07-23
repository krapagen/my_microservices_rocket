package config

type paymentClientConfig struct {
	Address string `yaml:"address" env:"PAYMENT_CLIENT_ADDRESS" env-default:"localhost:50052"`
}
