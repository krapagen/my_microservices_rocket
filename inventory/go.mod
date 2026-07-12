module github.com/krapagen/my_microservices_rocket/inventory

go 1.26.0

require (
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.79.2
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/brianvoe/gofakeit/v7 v7.15.0 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	go.opentelemetry.io/otel v1.42.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.42.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
)

replace github.com/krapagen/my_microservices_rocket/shared => ../shared
