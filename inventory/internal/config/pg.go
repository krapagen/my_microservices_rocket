package config

import "fmt"

type pgConfig struct {
	Host     string `yaml:"host"     env:"POSTGRES_HOST"     env-default:"localhost"`
	Port     string `yaml:"port"     env:"POSTGRES_PORT"     env-default:"5433"`
	Database string `yaml:"database" env:"POSTGRES_DB"       env-default:"inventory"`
	User     string `yaml:"user"     env:"POSTGRES_USER"     env-default:"inventory_admin"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env-default:"inventory_secret"`
	SSLMode  string `yaml:"sslmode"  env:"POSTGRES_SSLMODE"  env-default:"disable"`
}

func (c *pgConfig) DSN() string {
	// pool_max_conns: дефолт pgxpool — 4 соединения. Транзакции резервирования
	// ждут row-lock на горячих строках держа соединение — при маленьком пуле
	// ожидание блокирует даже чтения (ListParts) и запросы падают по deadline.
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s pool_max_conns=32",
		c.Host, c.Port, c.Database, c.User, c.Password, c.SSLMode,
	)
}
