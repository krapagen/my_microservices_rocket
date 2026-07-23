package app

// server, err := orderv1.NewServer(apiHandler, orderv1.WithErrorHandler(orderv1API.ErrorHandler))
import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	orderv1API "github.com/krapagen/my_microservices_rocket/order/internal/api/order/v1"
	inventoryClient "github.com/krapagen/my_microservices_rocket/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/krapagen/my_microservices_rocket/order/internal/client/grpc/payment/v1"
	orderRepository "github.com/krapagen/my_microservices_rocket/order/internal/repository/order"
	service "github.com/krapagen/my_microservices_rocket/order/internal/service/order"
	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

// NewHTTPHandler creates HTTP handler from gRPC clients (for tests)
func NewHTTPHandler(pool *pgxpool.Pool, txManager orderRepository.TxManager, inventoryGRPCClient inventoryv1.InventoryServiceClient, paymentGRPCClient paymentv1.PaymentServiceClient) (http.Handler, error) {
	// Repository layer
	orderRepo := orderRepository.New(pool, txManager)

	// Create client adapters
	invClient := inventoryClient.New(inventoryGRPCClient)
	payClient := paymentClient.New(paymentGRPCClient)

	// Service layer
	orderService := service.New(orderRepo, invClient, payClient, txManager)

	// API layer
	api := orderv1API.NewAPI(orderService)

	// Create OpenAPI server with error handler
	server, err := orderv1.NewServer(api, orderv1.WithErrorHandler(orderv1API.ErrorHandler))
	if err != nil {
		return nil, err
	}

	return server, nil
}
