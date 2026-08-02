package config

import "net"

type grpcConfig struct {
	Host string `yaml:"host" env:"GRPC_HOST" env-default:"0.0.0.0"`
	Port string `yaml:"port" env:"GRPC_PORT" env-default:"50053"`
}

func (c *grpcConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}
