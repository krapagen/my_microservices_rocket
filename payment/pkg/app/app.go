package app

import (
	"google.golang.org/grpc"

	paymentapi "github.com/krapagen/my_microservices_rocket/payment/internal/api/payment/v1"
	"github.com/krapagen/my_microservices_rocket/payment/internal/interceptor"
	service "github.com/krapagen/my_microservices_rocket/payment/internal/service/payment"
	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

// Interceptors возвращает grpc.ServerOption для тестов
func Interceptors() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	}
}

// RegisterServices регистрирует сервисы на gRPC сервере
func RegisterServices(grpcServer *grpc.Server) {
	svc := service.New()
	api := paymentapi.NewAPI(svc)
	paymentv1.RegisterPaymentServiceServer(grpcServer, api)
}
