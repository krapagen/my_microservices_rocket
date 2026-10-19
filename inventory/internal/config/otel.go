package config

// internal/config/otel.go
// otelConfig — внутренняя структура пакета config (имя с маленькой буквы),
// наружу торчит только метод OTel() OTelConfig интерфейса.
type otelConfig struct {
	Endpoint    string `yaml:"endpoint"`     // адрес OTEL Collector (OTLP gRPC)
	ServiceName string `yaml:"service_name"` // имя сервиса в трейсах / метриках / логах
}
