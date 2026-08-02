package config

type IAMClientConfig struct {
	Address string `yaml:"address"     env:"IAM_ADDRESS"     env-default:"localhost:50053"`
}

func (c *IAMClientConfig) Addr() string {
	return c.Address
}
