package config

type iamClientConfig struct {
	Address string `yaml:"address" env:"IAM_CLIENT_ADDRESS" env-default:"localhost:50053"`
}

func (c *iamClientConfig) Addr() string {
	return c.Address
}
